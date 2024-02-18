package database

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

const DB_TAG = "db"

// entity is a presentation of a entity that must have TableName function inside.
type entity interface {
	TableName() string
}

// FieldMap returns an list field names and a list pointer values of an entity.
func FieldMap[T entity](e T) ([]string, []any) {
	var fieldNames []string
	var fieldValues []any
	v := reflect.ValueOf(e).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Tag.Get(DB_TAG)
		fieldValue := v.Field(i).Addr().Interface()

		fieldNames = append(fieldNames, fieldName)
		fieldValues = append(fieldValues, fieldValue)
	}

	return fieldNames, fieldValues
}

// GetPlaceholders returns a string that grow from 1 to num with "$" in prefix and comma between them.
func GetPlaceholders(num int) string {
	result := []string{}
	for i := 1; i <= num; i++ {
		result = append(result, fmt.Sprintf("$%d", i))
	}

	return strings.Join(result, ", ")
}

// IsExistFieldInTable returns true if the field in params exists in entity.
func IsExistFieldInTable[T entity](dt T, target string) bool {
	t := reflect.TypeOf(dt)
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get(DB_TAG)
		fields = append(fields, tagValue)
	}

	return slices.Contains(fields, target)
}
