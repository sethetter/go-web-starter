package actions

import (
	"github.com/jmoiron/sqlx"
	"github.com/sethetter/go-web-starter/pkg/config"
	"github.com/sethetter/go-web-starter/pkg/services"
)

const (
	ErrInternalServerError = "Whoops! Something went wrong. Give it another try."
	ErrAuthCheckFailed     = "Sorry that code was incorrect, try again"
)

type Controller struct {
	DB           *sqlx.DB
	EmailService services.IEmailService
	Config       *config.Config
}
