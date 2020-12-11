package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var transformFunctionsMap = map[string]func(param string) string{
	"lower": func(value string) string {
		return strings.ToLower(value)
	},
	"upper": func(value string) string {
		return strings.ToUpper(value)
	},
}

func ReadJSON(reader io.Reader) interface{} {

	decoder := json.NewDecoder(reader)
	mapJSON := []map[string]interface{}{}
	for {
		var mapObject map[string]interface{}
		if err := decoder.Decode(&mapObject); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("%v\n", err)
			break

		}
		mapJSON = append(mapJSON, mapObject)
	}
	return mapJSON
}

func excuteFilters(filters, data map[string]interface{}, key string) {

	for action, value := range filters {
		switch action {
		case "newkey":
			{
				if val, ok := data[key]; ok {
					newKey := value.(string)
					delete(data, key)
					data[newKey] = val
					key = newKey

				}
			}
		case "transform":
			{
				if err := doTransform(data, key, transformFunctionsMap[value.(string)]); err != nil {
					fmt.Printf("%v\n", err)
				}

			}

		}
	}
}
func doTransform(data map[string]interface{}, key string, transform func(value string) string) error {

	if val, ok := data[key]; ok {
		data[key] = transform(val.(string))
	} else {
		return errors.New("key not found in json object - transform func faild")
	}
	return nil
}

func Filter(data, config map[string]interface{}) {

	for key, value := range data {
		if filters, ok := config[key]; ok {
			excuteFilters(filters.(map[string]interface{}), data, key)
		}
		switch value := value.(type) {
		case string, float64:
			break
		case []interface{}:
			for _, mapObject := range value {
				Filter(mapObject.(map[string]interface{}), config)
			}
		case interface{}:
			Filter(value.(map[string]interface{}), config)
		}
	}
}

func WriteJSON(data interface{}, writer io.Writer) {
	enc := json.NewEncoder(writer)
	if err := enc.Encode(&data); err != nil {
		fmt.Printf("%v\n", err)
	}
}
