package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"vilpanak/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitSession(cfg *config.Config) {
	Store = sessions.NewCookieStore([]byte(cfg.SessionSecret))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   cfg.Env == "production",
		SameSite: http.SameSiteLaxMode,
	}
}

// CSRF Protection Middleware
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := Store.Get(c.Request, "session_csrf")

		tokenVal, ok := session.Values["csrf_token"].(string)
		if !ok || tokenVal == "" {
			tokenBytes := make([]byte, 32)
			_, _ = rand.Read(tokenBytes)
			tokenVal = hex.EncodeToString(tokenBytes)
			session.Values["csrf_token"] = tokenVal
			_ = session.Save(c.Request, c.Writer)
		}

		c.Set("csrf_token", tokenVal)

		// Validate token for state-mutating requests
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
			submittedToken := c.PostForm("csrf_token")
			if submittedToken == "" {
				submittedToken = c.GetHeader("X-CSRF-Token")
			}

			if submittedToken == "" || submittedToken != tokenVal {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Invalid or missing CSRF token",
				})
				return
			}
		}

		c.Next()
	}
}

// RequireAuth ensures user is authenticated for admin endpoints
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := Store.Get(c.Request, "admin_session")
		if err != nil || session.Values["user_id"] == nil {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		c.Set("user_id", session.Values["user_id"])
		c.Set("user_name", session.Values["user_name"])
		c.Set("user_role", session.Values["user_role"])
		c.Next()
	}
}
