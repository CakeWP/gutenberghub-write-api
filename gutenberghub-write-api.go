package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pocketbase/pocketbase"

	"github.com/pocketbase/pocketbase/core"
)

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
