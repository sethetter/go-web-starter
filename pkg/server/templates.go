package server

import (
	"html/template"
	"path"

	"github.com/gin-contrib/multitemplate"
)

func renderer(templatePath string) multitemplate.Renderer {
	funcMap := template.FuncMap{
		"formatAsDate":          formatAsDate,
		"formatAsRfc3339String": formatAsRfc3339String,
	}

	r := multitemplate.NewRenderer()

	// Pages
	basePath := path.Join(templatePath, "pages/base.html")
	r.AddFromFilesFuncs("page.index", funcMap, basePath, path.Join(templatePath, "pages/index.html"))
	r.AddFromFilesFuncs("page.login", funcMap, basePath, path.Join(templatePath, "pages/login.html"))

	// Partials
	r.AddFromFilesFuncs("partials.login-start", funcMap, path.Join(templatePath, "partials/login-start.html"))
	r.AddFromFilesFuncs("partials.login-verify", funcMap, path.Join(templatePath, "partials/login-verify.html"))

	return r
}
