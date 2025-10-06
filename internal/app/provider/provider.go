package provider

import (
	"context"
	"fmt"

	"github.com/ttrtcixy/workout/internal/app/closer"
	"github.com/ttrtcixy/workout/internal/config"
	"github.com/ttrtcixy/workout/internal/core/repository"
	"github.com/ttrtcixy/workout/internal/core/usecase"
	"github.com/ttrtcixy/workout/internal/delivery/grpc"
	"github.com/ttrtcixy/workout/internal/delivery/http"
	"github.com/ttrtcixy/workout/internal/delivery/http/handlers"
	apperrors "github.com/ttrtcixy/workout/internal/errors"
	"github.com/ttrtcixy/workout/internal/logger"
	storage "github.com/ttrtcixy/workout/internal/storage/pg"
)

type Provider struct {
	logger logger.Logger
	closer closer.Closer

	cfg *config.Config

	db storage.DB

	authClient *grpc.AuthClient

	handlers   *handlers.Handlers
	usecase    *usecase.UseCase
	repository *repository.Repository

	httpServer *http.Server
}

func New(ctx context.Context) (p *Provider, err error) {
	const op = "Provider.New"

	p = &Provider{}

	if err = p.initLogger(); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initConfig(); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initCloser(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	defer func() {
		if err != nil {
			p.closer.Close()
		}
	}()

	p.closer.Add(
		"env clear",
		p.cfg.Close,
	)

	if err = p.initDB(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initServices(); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initAuthClient(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initRepository(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initUseCase(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initHandler(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}

	if err = p.initHTTPServer(ctx); err != nil {
		return p, apperrors.Wrap(op, err)
	}
	p.closer.Add("stop http server", p.httpServer.Close)

	return p, nil
}

func (p *Provider) initLogger() error {
	const op = "Provider.initLogger"
	p.logger = logger.Load()
	return nil
}

func (p *Provider) initConfig() error {
	const op = "Provider.initConfig"
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("op: %s - config init failed: %w", op, err)
	}

	p.cfg = cfg
	p.logger.Info("[+] config loaded")
	return nil
}

func (p *Provider) initCloser(ctx context.Context) error {
	const op = "Provider.initCloser"
	p.closer = closer.New(closer.Config{
		TotalDuration: p.cfg.Closer.TotalDuration(),
		FuncDuration:  p.cfg.Closer.FuncDuration(),
		Logger:        p.logger,
	})

	p.logger.Info("[+] closer loaded")
	return nil
}

func (p *Provider) initDB(ctx context.Context) error {
	const op = "Provider.initDB"
	db, err := storage.New(ctx, p.logger, p.cfg.DB)
	if err != nil {
		return fmt.Errorf("op: %s - db init failed: %w", op, err)
	}
	p.db = db
	p.logger.Info("[+] connect to database successful")

	p.closer.Add(
		"close db connection",
		p.db.Close,
	)

	return nil
}

func (p *Provider) initServices() error {
	const op = "Provider.initServices"

	return nil
}

func (p *Provider) initAuthClient(ctx context.Context) (err error) {
	const op = "Provider.initAuthClient"

	if p.authClient, err = grpc.NewAuthClient(p.logger, p.cfg.GrpcAuthServer); err != nil {
		return fmt.Errorf("op: %s - auth client init failed: %w", op, err)
	}

	p.closer.Add(
		"close grpc auth client connection",
		p.authClient.Close,
	)

	return nil
}

func (p *Provider) initRepository(ctx context.Context) error {
	const op = "Provider.initRepository"
	p.repository = repository.NewRepository(ctx, p.logger, p.db)

	return nil
}

func (p *Provider) initUseCase(ctx context.Context) error {
	const op = "Provider.initUseCase"
	p.usecase = usecase.NewUseCase(p.logger, p.repository)
	return nil
}

func (p *Provider) initHandler(ctx context.Context) error {
	const op = "Provider.initHandler"
	p.handlers = handlers.NewHandlers(p.logger, p.usecase)
	return nil
}

func (p *Provider) initHTTPServer(ctx context.Context) error {
	const op = "Provider.initHTTPServer"
	p.httpServer = http.New(p.cfg.HttpServer, p.logger, p.handlers)
	return nil
}

func (p *Provider) HTTPServer() *http.Server {
	return p.httpServer
}

func (p *Provider) Logger() logger.Logger {
	return p.logger
}

func (p *Provider) Closer() closer.Closer {
	return p.closer
}
