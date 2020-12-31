package middleware

import (
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/representations"
)

const identityKey = "id"

func NewAuthMiddleware() *jwt.GinJWTMiddleware {
	adminUsername := "admin"
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:       "opstic",
		Key:         []byte("sa76duh387dfsihuasdf897ui398dfsuio"),
		Timeout:     365 * 24 * time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entities.User); ok {
				// log.Printf("payload user %s", v.Username)
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			// log.Printf("identity %s", claims[identityKey].(string))
			return &entities.User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginValues representations.Login
			if err := c.Bind(&loginValues); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginValues.Username
			password := loginValues.Password
			// TODO: Use database?
			// log.Printf("authenticate %s %s", username, password)

			if os.Getenv("DEVOPSTIC_USERNAME") != "" {
				adminUsername = os.Getenv("DEVOPSTIC_USERNAME")
			}
			adminPassword := "admin"
			if os.Getenv("DEVOPSTIC_PASSWORD") != "" {
				adminPassword = os.Getenv("DEVOPSTIC_PASSWORD")
			}
			if username == adminUsername && password == adminPassword {
				return &entities.User{
					Username:  username,
					LastName:  "Devopstic",
					FirstName: "Admin",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// log.Printf("authorize %s", data)
			if v, ok := data.(*entities.User); ok && v.Username == adminUsername {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt, param: token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:         time.Now,
		SigningAlgorithm: "HS256",
	}
	authMiddleware, err := jwt.New(authMiddleware)
	if err != nil {
		log.Fatal(err)
	}
	return authMiddleware
}
