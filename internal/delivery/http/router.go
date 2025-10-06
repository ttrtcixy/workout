package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ttrtcixy/workout/internal/delivery/http/handlers"
	"github.com/ttrtcixy/workout/internal/delivery/http/middlewares"
	"net/http"
)

type Router struct {
	router *gin.Engine
}

// Handler return http.Handler
func (r *Router) Handler() http.Handler {
	return r.router.Handler()
}

// NewRouter create Router, init Routes and Middlewares
func NewRouter(handlers *handlers.Handlers) *Router {
	r := &Router{router: gin.Default()}

	r.initRoutes(handlers)
	r.initMiddleware()

	return r
}

// initRoutes add handlers to routes
func (r *Router) initRoutes(handlers *handlers.Handlers) {
	r.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.router.POST("/create-workout", handlers.CreateWorkout.Run)
}

// initMiddleware add middlewares to routes
func (r *Router) initMiddleware() {
	r.router.Use(middlewares.ShutdownMiddleware())
}
