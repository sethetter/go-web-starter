package actions

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sethetter/go-web-starter/pkg/middleware"
	"github.com/sethetter/go-web-starter/pkg/models"
)

func (ctrl *Controller) LoginPage(ctx *gin.Context) {
	ctx.HTML(200, "page.login", addGlobalVars(ctx, gin.H{}))
}

func (ctrl *Controller) LoginStartForm(ctx *gin.Context) {
	ctx.HTML(200, "partials.login-start", addGlobalVars(ctx, gin.H{}))
}

func (ctrl *Controller) LoginStart(ctx *gin.Context) {
	email := ctx.PostForm("email")
	newUser := models.NewUser{Email: email}

	if valid, errs := newUser.Validate(); !valid {
		ctx.HTML(200, "partials.login-start", gin.H{
			"formErrors": formErrorsFromValidation(govalidator.Errors(errs)),
		})
		return
	}

	user, err := newUser.InsertOrGet(ctrl.DB)
	if err != nil {
		log.Printf("newUser.InsertOrGet failed: %s", err)
		serverError(ctx, "partials.login-start", gin.H{"email": email})
		return
	}

	challenge, err := models.NewAuthChallenge(user, ctrl.DB)
	if err != nil {
		log.Printf("models.NewAuthChallenge failed: %s", err)
		serverError(ctx, "partials.login-start", gin.H{"email": email})
		return
	}

	message := fmt.Sprintf("Your login code is: %s", challenge.Code)
	ctrl.EmailService.SendEmail(ctrl.Config.AppName, user.Email, "Verify your login", message)

	session := sessions.Default(ctx)
	session.Set(middleware.SessKey_User, user.ID)
	session.Set(middleware.SessKey_Authenticated, false)
	session.Save()

	ctx.HTML(200, "partials.login-verify", addGlobalVars(ctx, gin.H{"email": email}))
}

func (ctrl *Controller) LoginVerify(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userId := session.Get(middleware.SessKey_User).(string)
	code := ctx.PostForm("code")

	challenge, err := models.AuthChallengeForUserID(userId, ctrl.DB)
	if err != nil {
		// If we didn't find a challenge, they need to restart the flow
		if err == sql.ErrNoRows {
			serverError(ctx, "partials.login-start", gin.H{})
			return
		}
		serverError(ctx, "partials.login-verify", gin.H{"code": code})
		return
	}

	if !challenge.Check(code) {
		ctx.HTML(200, "partials.login-start", gin.H{
			"errors": [][]string{{"warning", ErrAuthCheckFailed}},
		})
		return
	}

	session.Set(middleware.SessKey_Authenticated, true)
	session.AddFlash("Logged in successfully!", "success")
	session.Save()

	ctx.Header("hx-redirect", "/")
}

func (ctrl *Controller) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	session.Set(middleware.SessKey_User, nil)
	session.Set(middleware.SessKey_Authenticated, false)

	session.AddFlash("Logged out successfully", "success")

	session.Save()

	ctx.Redirect(301, "/")
}
