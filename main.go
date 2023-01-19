package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// RateLimitCollectionOperations is a middleware that rate limits HTTP requests on collection operations.
// Possible operations:
// - list
// - view
// - create
// - update
// - delete
func RateLimitCollectionOperations(collection string, operations []string, limiter echo.MiddlewareFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contains := func(s string, list []string) bool {
				for _, v := range list {
					if v == s {
						return true
					}
				}
				return false
			}

			collectionEndpoint := fmt.Sprintf("/api/collections/%s/records", collection)

			if contains("list", operations) && c.Request().URL.Path == collectionEndpoint && c.Request().Method == http.MethodGet {
				return limiter(next)(c)
			} else if contains("view", operations) && strings.HasPrefix(c.Request().URL.Path, collectionEndpoint) && c.Request().Method == http.MethodGet {
				return limiter(next)(c)
			} else if contains("create", operations) && strings.HasPrefix(c.Request().URL.Path, collectionEndpoint) && c.Request().Method == http.MethodPost {
				return limiter(next)(c)
			} else if contains("update", operations) && strings.HasPrefix(c.Request().URL.Path, collectionEndpoint) && c.Request().Method == http.MethodPatch {
				return limiter(next)(c)
			} else if contains("delete", operations) && strings.HasPrefix(c.Request().URL.Path, collectionEndpoint) && c.Request().Method == http.MethodDelete {
				return limiter(next)(c)
			}

			return next(c)
		}
	}
}

func main() {
	app := pocketbase.New()

	// Rate limiter.
	limiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2))
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Rate limit list and view operations on the "posts" collection
		e.Router.Use(RateLimitCollectionOperations("posts", []string{"list", "view"}, limiter))
		return nil
	})

	// Adding a new secured `Connect Project` Route.
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/connect",
			Handler: func(c echo.Context) error {

				// Project Access Key.
				// accessKey := c.Request().Header["gutenberghub-access-key"]
				accessKey := c.QueryParam("key")

				return c.String(http.StatusOK, "Your Access Key: "+accessKey)
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
