package grpc

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/workout/internal/config"
	apperrors "github.com/ttrtcixy/workout/internal/errors"
	"github.com/ttrtcixy/workout/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client users.UsersAuthClient
	log    logger.Logger
	cfg    *config.GRPCAuthServer

	reconnecting atomic.Bool
	mu           sync.Mutex
}

var ErrClientNotInitialized = errors.New("client not initialized")
var ErrServerUnavailable = errors.New("grpc auth server unavailable")

type AuthClientRequest struct {
	AccessToken  string
	RefreshToken *string
}

func NewAuthClient(log logger.Logger, cfg *config.GRPCAuthServer) (*AuthClient, error) {
	const op = "grpc.NewAuthClient"
	client := &AuthClient{
		log: log,
		cfg: cfg,
	}

	if err := client.connect(); err != nil {
		return nil, apperrors.Wrap(op, err)
	}

	return client, nil
}

func (a *AuthClient) connect() (err error) {
	const op = "AuthClient.connect"
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.conn != nil {
		if closeErr := a.conn.Close(); closeErr != nil {
			a.log.Error("%s: closing connection err: %w", op, err)
		}
	}

	// todo use config
	opts := []grpc.DialOption{
		// security
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// connection params
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  3 * time.Second,
				Multiplier: 1.6,
				MaxDelay:   30 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
		// keepalive params
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 5 * time.Second}),
		// grpc request retry
		grpc.WithDefaultServiceConfig(`{
			"retryPolicy": {
				"maxAttempts": 4,       
				"initialBackoff": "0.1s", 
				"maxBackoff": "1s",    
				"backoffMultiplier": 2,
				"retryableStatusCodes": ["UNAVAILABLE"]
			}
		}`),
	}

	// todo add ping
	if a.conn, err = grpc.NewClient(a.cfg.Addr(), opts...); err != nil {
		return apperrors.Wrap(op, err)
	}

	a.client = users.NewUsersAuthClient(a.conn)
	return nil
}

func (a *AuthClient) getClient() (users.UsersAuthClient, error) {
	const op = "AuthClient.getClient"

	if a.client == nil {
		return nil, apperrors.Wrap(op, ErrClientNotInitialized)
	}

	if a.conn.GetState() == connectivity.Shutdown {
		if a.reconnecting.CompareAndSwap(false, true) {
			go a.reconnect(context.Background())
		}
		return nil, apperrors.Wrap(op, ErrServerUnavailable)
	}

	return a.client, nil
}

func (a *AuthClient) reconnect(ctx context.Context) {
	const op = "AuthClient.reconnect"
	defer a.reconnecting.Store(false)

	a.log.Error("%s: grpc auth server is not available, reconnection started", op)
	var recAttempts int
	for {
		select {
		case <-ctx.Done():
			a.log.Info("%s: server reconnection stopped", op)
			return
		default:
			if err := a.connect(); err != nil {
				recAttempts++
				a.log.Error("%s: grpc auth server is not available, reconnection attempt %d, err: %w", op, recAttempts, err)
				continue
			}

			if a.conn.GetState() == connectivity.Ready {
				a.log.Info("%s: grpc auth server reconnected, reconnection attempt %d", op, recAttempts)
				return
			}
		}
	}
}

func (a *AuthClient) VerifyUser(ctx context.Context, payload *AuthClientRequest) (response *users.VerifyAccessTokenResponse, err error) {
	const op = "AuthClient.ValidateUserAuth"
	defer func() {
		if err != nil {
			var userErr apperrors.UserError
			if errors.As(err, &userErr) {
				return
			}
			a.log.ErrorOp(op, err)
			err = apperrors.ErrServer
		}
	}()
	client, err := a.getClient()
	if err != nil {
		return nil, err
	}

	response, err = client.VerifyToken(ctx, &users.VerifyAccessTokenRequest{
		AccessToken:  payload.AccessToken,
		RefreshToken: payload.RefreshToken,
	})
	if err != nil {
		// todo log auth server errors, and return user error
		return nil, err
	}

	return response, nil
}

// todo test with error
// todo могут быть ошибки если начнется реконект а мы закрываем соединение
func (a *AuthClient) Close(ctx context.Context) error {
	const op = "AuthClient.close"

	if err := a.conn.Close(); err != nil {
		return apperrors.Wrap(op, err)
	}

	return nil
}
