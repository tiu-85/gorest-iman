package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func ReflectID(model interface{}) reflect.Value {
	modelVal := reflect.ValueOf(model).Elem()
	idVal := modelVal.FieldByNameFunc(func(s string) bool {
		return strings.ToLower(s) == "id"
	})
	return idVal
}

func ReflectHasID(model interface{}) (bool, error) {
	idVal := ReflectID(model)
	switch idVal.Kind() {
	case reflect.String:
		if idVal.String() == "" {
			return false, nil
		}
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16:
		if idVal.Int() == 0 {
			return false, nil
		}
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16:
		if idVal.Uint() == 0 {
			return false, nil
		}
	default:
		return false, fmt.Errorf("can`t detect ID field")
	}
	return true, nil
}
