package devopstic

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"github.com/splisson/devopstic/handlers"
	"github.com/splisson/devopstic/middleware"
	"os"
)

func BuildEngine(commitHandlers *handlers.CommitHandlers, eventHandlers *handlers.EventHandlers, githubEventHandlers *handlers.GithubEventHandlers, gitlabEventHandlers *handlers.GitlabEventHandlers, pagerDutyHandlers *handlers.PagerDutyHandlers) *gin.Engine {

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

	githubSecret := os.Getenv("DEVOPSTIC_GITHUB_WEBHOOK_SECRET")
	if len(githubSecret) <= 0 {
		log.Warn("missing github webhook secret env var: DEVOPSTIC_GITHUB_WEBHOOK_SECRET")
	} else {
		githubAuthMiddleware := middleware.NewGithubAuthMiddleware(githubSecret)
		githubAuth := r.Group("/github")
		githubAuth.Use(githubAuthMiddleware.MiddlewareFunc())
		{
			githubAuth.POST("/events", githubEventHandlers.PostGithubEvents)
		}
	}
	gitlabSecret := os.Getenv("DEVOPSTIC_GITLAB_WEBHOOK_SECRET")
	if len(gitlabSecret) <= 0 {
		log.Warn("missing gitlab webhook secret env var: DEVOPSTIC_GITLAB_WEBHOOK_SECRET")
	} else {
		gitlabAuthMiddleware := middleware.NewGitlabAuthMiddleware(gitlabSecret)
		gitlabAuth := r.Group("/gitlab")
		gitlabAuth.Use(gitlabAuthMiddleware.MiddlewareFunc())
		{
			gitlabAuth.POST("/events", gitlabEventHandlers.PostGitlabEvents)
		}
	}

	// PagerDuty
	pagerDutyAuth := r.Group("/pagerduty")
	pagerDutyAuth.POST("/incidents", pagerDutyHandlers.PostPagerDutyIncidents)

	return r
}
