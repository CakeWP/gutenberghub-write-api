package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/pocketbase/pocketbase"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

func excludeFields(data interface{}, fields []string) ([]byte, error) {
	// check if data is of type array
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return nil, fmt.Errorf("data should be of type slice, got %T", data)
	}

	// convert data to json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Error while marshalling data: %v", err)
	}

	// Unmarshal the JSON data into a map
	var dataMap []map[string]interface{}
	err = json.Unmarshal(jsonData, &dataMap)
	if err != nil {
		return nil, fmt.Errorf("Error while unmarshalling data: %v", err)
	}

	// Iterate through the map and delete the specified fields
	for _, item := range dataMap {
		for _, field := range fields {
			delete(item, field)
		}
	}

	// Marshal the modified data back to JSON
	jsonData, err = json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("Error while marshalling data: %v", err)
	}

	return jsonData, nil
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

	// Adding a new query field parameter when listing view.
	app.OnRecordsListRequest().Add(func(e *core.RecordsListEvent) error {

		hasExclusion := e.HttpContext.Request().URL.Query().Has("excluded")

		if !hasExclusion {
			return nil
		}

		excludedFields := strings.Split(e.HttpContext.Request().URL.Query().Get("excluded"), ",")

		items := e.Result.Items
		excludedItems, err := excludeFields(items, excludedFields)

		if err != nil {
			fmt.Println(err)
		}

		var updatedItems []map[string]interface{}

		json.Unmarshal(excludedItems, &updatedItems)

		e.Result.Items = updatedItems

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
