package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/kimihito-sandbox/pocketbase-templ-esbuild/templates"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

var (
	//go:embed assets/dist/*
	embededAssets embed.FS

	distDirFS = echo.MustSubFS(embededAssets, "assets")
)

func homeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Home())
}

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Pre(middleware.RemoveTrailingSlash())

		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		e.Router.GET("/assets/*", apis.StaticDirectoryHandler(distDirFS, false))
		e.Router.GET("/hello", homeHandler)
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
