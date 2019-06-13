package opstic

import (
	"github.com/gin-gonic/gin"
	"github.com/splisson/opstic/handlers"
	"github.com/splisson/opstic/middleware"
)

func BuildEngine(eventHandlers *handlers.EventHandlers) *gin.Engine {

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

	r.POST("/login", authMiddleware.LoginHandler)

	r.GET("/events", eventHandlers.GetEvents)

	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/events", eventHandlers.PostEvents)
		auth.POST("/webhook/:token/events", eventHandlers.PostEvents)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
		auth.GET("/hello", handlers.HelloHandler)
	}

	return r
}
