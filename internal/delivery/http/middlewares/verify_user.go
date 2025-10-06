package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ttrtcixy/workout/internal/delivery/grpc"
	"github.com/ttrtcixy/workout/internal/logger"
)

const UserInfoKey = "UserInfo"

type AuthMiddleware struct {
	authClient *grpc.AuthClient
	log        logger.Logger
}

func NewAuthMiddleware(authClient *grpc.AuthClient, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient, log: log}
}

func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Request.Cookie("access_token")
		if err != nil {
			fmt.Println(err)
			return
		}

		// todo check if token exp, if true add refresh token to request, else only access token

		ctx := context.WithValue(c.Request.Context(), UserInfoKey)

		c.Request.WithContext()
		c.Next()
	}
}
