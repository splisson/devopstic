package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GithubAuthMiddleware struct {
	webhookSecret string
}

var (
	ErrMissingSignatureHeader = errors.New("missing header X-Hub-Signature")
)

func NewGithubAuthMiddleware(webhookSecret string) *GithubAuthMiddleware {
	mw := new(GithubAuthMiddleware)
	mw.webhookSecret = webhookSecret
	return mw
}

func (mw *GithubAuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

func (mw *GithubAuthMiddleware) unauthorized(c *gin.Context, code int, message string) {
	//if !mw.DisabledAbort {
	//    c.Abort()
	//}
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (mw *GithubAuthMiddleware) middlewareImpl(c *gin.Context) {
	// Get receivedSignature from header
	receivedSignature := c.GetHeader("X-Hub-Signature")

	if len(receivedSignature) <= 0 {
		mw.unauthorized(c, http.StatusBadRequest, ErrMissingSignatureHeader.Error())
		return
	}

	if !isValidSignature(c.Request, mw.webhookSecret) {
		mw.unauthorized(c, http.StatusUnauthorized, "invalid signature")
		return
	}

	c.Set("id", "admin") // TODO: dynamic user

	c.Next()
}

func isValidSignature(request *http.Request, key string) bool {
	// Assuming a non-empty header
	gotHash := strings.SplitN(request.Header.Get("X-Hub-Signature"), "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}

	body, err := ioutil.ReadAll(request.Body)
	request.Body.Close()
	if err != nil {
		log.Printf("cannot read the request body: %s\n", err)
		return false
	}

	hash := hmac.New(sha1.New, []byte(key))
	if _, err := hash.Write(body); err != nil {
		log.Printf("cannot compute the HMAC for request: %s\n", err)
		return false
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	log.Println("EXPECTED HASH:", expectedHash)

	// Restore request body
	request.Body = ioutil.NopCloser(bytes.NewReader(body))

	return gotHash[1] == expectedHash
}
