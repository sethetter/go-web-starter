package middleware

import (
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sethetter/go-web-starter/pkg/config"
)

func SecurityHeaders(config *config.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if config.Env != "debug" {
			ctx.Header("strict-transport-security", "max-age=31536000; includeSubdomains")
		}

		ctx.Header("content-security-policy", "default-src 'self';")
		ctx.Header("x-content-type-options", "nosniff")
		ctx.Header("x-frame-options", "DENY")
		ctx.Header("content-origin-resource-policy", "same-origin")
		ctx.Header("cross-origin-opener-policy", "same-origin-allow-popups")
		ctx.Header("cross-origin-embedder-policy", "require-corp")

		ctx.Next()
	}
}

func Csrf(secret string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Make sure we set before checking omitted routes, so a token
		// is present on subsequent requests
		token := setOrGetCsrfToken(session, secret)

		if ctx.Request.Method == "GET" {
			ctx.Next()
			return
		}

		// Inspect request content-type to determine
		// where we get the token to validate from
		switch strings.ToLower(ctx.GetHeader("content-type")) {
		case "application/x-www-form-urlencoded":
			if ctx.PostForm("_csrf") == token {
				ctx.Next()
				return
			}
			log.Printf("middleware.Csrf: mismatch on csrf token: %s", ctx.PostForm("_csrf"))
		default:
			log.Printf("middleware.Csrf: unhandled content-type: %s", ctx.GetHeader("content-type"))
		}
		// TODO: fail with a specific CSRF error (AbortWithError)?
		ctx.AbortWithStatus(403)
	}
}

func setOrGetCsrfToken(session sessions.Session, secret string) string {
	csrf := session.Get(SessKey_Csrf)
	if csrf, ok := csrf.(string); ok && csrf != "" {
		return csrf
	}

	// Generate and save a new token
	token := uuid.NewString()

	session.Set(SessKey_Csrf, token)
	session.Save()

	return token
}
