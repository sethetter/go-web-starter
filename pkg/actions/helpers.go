package actions

import (
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sethetter/go-web-starter/pkg/middleware"
)

func addGlobalVars(ctx *gin.Context, base gin.H) gin.H {
	session := sessions.Default(ctx)

	flashMap := map[string][]interface{}{}
	flashMap["danger"] = session.Flashes("danger")
	flashMap["warning"] = session.Flashes("warning")
	flashMap["success"] = session.Flashes("success")
	flashMap["info"] = session.Flashes("info")
	flashMap["other"] = session.Flashes()
	base["flashes"] = flashMap

	base["csrf"] = session.Get(middleware.SessKey_Csrf)

	session.Save()

	return base
}

func formErrorsFromValidation(errs []error) map[string][]string {
	out := make(map[string][]string)
	out["rest"] = []string{}

	for _, err := range errs {
		parts := strings.SplitN(err.Error(), ": ", 2)
		if len(parts) != 2 {
			out["rest"] = append(out["rest"], err.Error())
		}

		key, str := parts[0], parts[1]

		if v, ok := out[key]; ok {
			out[key] = append(v, str)
		} else {
			out[key] = []string{str}
		}
	}

	return out
}

func serverError(ctx *gin.Context, partial string, inVars gin.H) {
	vars := gin.H{}

	for k, v := range inVars {
		vars[k] = v
	}

	vars["errors"] = [][]string{{"danger", ErrInternalServerError}}

	ctx.HTML(200, partial, addGlobalVars(ctx, vars))
}
