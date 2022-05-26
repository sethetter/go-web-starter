package actions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sethetter/go-web-starter/pkg/middleware"
)

func (ctrl *Controller) Index(ctx *gin.Context) {
	session := sessions.Default(ctx)
	vars := gin.H{}

	authenticated := session.Get("authenticated")
	if a, ok := authenticated.(bool); a && ok {
		vars["user"] = session.Get(middleware.SessKey_User)
	}

	ctx.HTML(200, "page.index", addGlobalVars(ctx, vars))
}
