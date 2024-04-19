package http

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/collector"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/resources/views"
	"github.com/labstack/echo/v4"
)

func NewStatusPageHandler(app *app.App, cache *ResultCache) echo.HandlerFunc {
	tmpl := template.Must(template.ParseFS(views.Views, "index.html"))

	return func(c echo.Context) error {
		var buf bytes.Buffer
		data := struct {
			Config  config.Config
			Results []collector.Result
		}{
			Config:  *app.Config,
			Results: cache.Get(),
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			log.Printf("ERROR - could not render template: %s", err)
			return err
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}
