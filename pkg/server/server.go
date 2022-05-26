package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sethetter/go-web-starter/pkg/actions"
	"github.com/sethetter/go-web-starter/pkg/config"
	"github.com/sethetter/go-web-starter/pkg/middleware"
	"github.com/sethetter/go-web-starter/pkg/services"
)

type ServerConfig struct {
	Config       *config.Config
	DB           *sql.DB
	EmailService services.IEmailService
	TemplatePath string
}

func NewServer(c *ServerConfig) (http.Server, error) {
	gin.SetMode(c.Config.Env)
	gin.DefaultWriter = log.Writer()

	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		return http.Server{}, fmt.Errorf("failed to SetTrustedProxies: %w", err)
	}

	sessionOpts := sessions.Options{
		Path:     "/",
		MaxAge:   24 * 60, // 1 day
		Secure:   c.Config.Env != "debug",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	sessionStore, err := postgres.NewStore(c.DB, []byte(c.Config.AppSecret))
	if err != nil {
		return http.Server{}, fmt.Errorf("failed to open postgres session store: %w", err)
	}

	sessionStore.Options(sessionOpts)
	router.Use(sessions.Sessions("mysession", sessionStore))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{c.Config.URL}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.Csrf(c.Config.AppSecret))
	router.Use(middleware.SecurityHeaders(c.Config))

	sqlxDb := sqlx.NewDb(c.DB, "postgres")

	ctrl := &actions.Controller{
		DB:           sqlxDb,
		Config:       c.Config,
		EmailService: c.EmailService,
	}

	router.HTMLRender = renderer(c.TemplatePath)
	registerRoutes(ctrl, router)

	return http.Server{
		Addr:    c.Config.Port,
		Handler: router,
	}, nil
}
