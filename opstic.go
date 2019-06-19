package devopstic

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/devopstic/handlers"
	"github.com/splisson/devopstic/middleware"
)

func BuildEngine(commitHandlers *handlers.CommitHandlers, eventHandlers *handlers.EventHandlers) *gin.Engine {

	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// JWT
	// the jwt middleware
	authMiddleware := middleware.NewAuthMiddleware()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	r.POST("/tokens", authMiddleware.LoginHandler)

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/commits", commitHandlers.GetCommits)
		auth.GET("/events", eventHandlers.GetEvents)
		auth.POST("/events", eventHandlers.PostEvents)
		auth.POST("/webhook/:token/events", eventHandlers.PostEvents)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	return r
}
