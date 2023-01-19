package main

import (
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// Route: Connect Project Route
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/connect",
			Handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "I'm zafar kamal!!")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(app),
				// apis.RequireAdminAuth(),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
