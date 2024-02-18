package reflect_util

import (
	"encoding/json"
	"reflect"
	"strconv"
)

// ConvertMapToStruct convert the map[string]any to a passing result of generic struct.
// TODO: recursive all map and convert to child struct
func ConvertMapToStruct[T any](m map[string]any, s *T) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		name := field.Tag.Get("json")
		if value, ok := m[name]; ok {
			switch {
			// convert string to int64 if the result struct defined int64.
			case reflect.TypeOf(value).Kind() == reflect.String && field.Type.Kind() == reflect.Int64:
				iVal, err := strconv.ParseInt(value.(string), 10, 64)
				if err != nil {
					return err
				}
				stValue.Field(i).Set(reflect.ValueOf(iVal))
				// TODO: convert other concrete types.
			case reflect.TypeOf(value).Kind() == reflect.Float64 && field.Type.Kind() == reflect.Int64:
				stValue.Field(i).Set(reflect.ValueOf(int64(value.(float64))))
			default:
				stValue.Field(i).Set(reflect.ValueOf(value))
			}
		}
	}

	return nil
}

// CopyStruct copy value of source to destination by using json marshal and unmarshal json methods.
func CopyStruct[T, V any](source *V, destination *T) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, destination); err != nil {
		return err
	}

	return nil
}
