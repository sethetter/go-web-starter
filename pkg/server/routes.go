package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sethetter/go-web-starter/pkg/actions"
)

func registerRoutes(ctrl *actions.Controller, router *gin.Engine) {
	router.Static("/assets", "assets")
	router.GET("/", ctrl.Index)
	router.GET("/login", ctrl.LoginPage)
	router.GET("/logout", ctrl.Logout)
	router.GET("/login/start", ctrl.LoginStartForm)
	router.POST("/login/start", ctrl.LoginStart)
	router.POST("/login/verify", ctrl.LoginVerify)
}
