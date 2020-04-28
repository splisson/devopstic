package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GitlabAuthMiddleware struct {
	webhookSecret string
}

var (
	ErrMissingTokenHeader = errors.New("missing header X-Gitlab-Token")
)

func NewGitlabAuthMiddleware(webhookSecret string) *GitlabAuthMiddleware {
	mw := new(GitlabAuthMiddleware)
	mw.webhookSecret = webhookSecret
	return mw
}

func (mw *GitlabAuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

func (mw *GitlabAuthMiddleware) unauthorized(c *gin.Context, code int, message string) {
	//if !mw.DisabledAbort {
	//    c.Abort()
	//}
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (mw *GitlabAuthMiddleware) middlewareImpl(c *gin.Context) {
	// Get receivedToken from header
	receivedToken := c.GetHeader("X-Gitlab-Token")

	if len(receivedToken) <= 0 {
		mw.unauthorized(c, http.StatusBadRequest, ErrMissingTokenHeader.Error())
		return
	}

	if !isValidToken(c.Request, mw.webhookSecret) {
		mw.unauthorized(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set("id", "admin") // TODO: dynamic user

	c.Next()
}

func isValidToken(request *http.Request, expectedToken string) bool {
	// Assuming a non-empty header
	token := request.Header.Get("X-Gitlab-Token")
	if len(token) > 0 && token == expectedToken {
		return true
	}
	return false
}
