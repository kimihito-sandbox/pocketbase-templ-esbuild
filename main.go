package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/kimihito-sandbox/pocketbase-templ-esbuild/templates"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

//go:embed assets/dist/*
var embededAssets embed.FS

func main() {
	app := pocketbase.New()
	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
	e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
	e.Router.GET("/assets/*", apis.StaticDirectoryHandler(echo.MustSubFS(embededAssets, "assets"), false))
	e.Router.GET("/home", func(c echo.Context) error {
		return Render(c, http.StatusOK, templates.Home())
	})
	e.Router.GET("/hello/:name", func(c echo.Context) error {
		name := c.PathParam("name")
		return c.JSON(http.StatusOK, "Hello, "+name+"!")
	})
		return nil
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
