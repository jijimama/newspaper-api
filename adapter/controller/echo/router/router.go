package router

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"newspaper-api/adapter/controller/echo/handler"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewEchoRouter(db *gorm.DB) *echo.Echo {
	router := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./adapter/presenter/html/*")),
	}
	router.Renderer = renderer
	router.GET("/", handler.Index)
	router.GET("/health", handler.Health)
	return router
}
