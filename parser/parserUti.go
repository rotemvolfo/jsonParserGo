package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var transformFunctionsMap = map[string]func(param ...string) (error, string){

	"lower": func(param ...string) (error, string) {
		return nil, strings.ToLower(param[0])
	},
	"upper": func(param ...string) (error, string) {
		return nil, strings.ToUpper(param[0])
	},
}

func ReadJSON(reader io.Reader) ([]map[string]interface{}, error) {

	decoder := json.NewDecoder(reader)
	arrJSON := []map[string]interface{}{}
	for {
		var mapObject map[string]interface{}
		if err := decoder.Decode(&mapObject); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		arrJSON = append(arrJSON, mapObject)
	}
	return arrJSON, nil
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
func doTransform(data map[string]interface{}, key string, transform func(param ...string) (error, string)) error {
	var err error
	if val, ok := data[key]; ok {
		if err, data[key] = transform(val.(string)); err != nil {
			return err
		}
		return nil
	}
	return errors.New("key not found in json object - transform func failed")
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

func WriteJSON(data interface{}, writer io.Writer) error {
	enc := json.NewEncoder(writer)
	if err := enc.Encode(&data); err != nil {
		return err
	}
	return nil
}
