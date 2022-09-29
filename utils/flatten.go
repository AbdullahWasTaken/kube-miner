package utils

import (
	"reflect"
	"strconv"
)

func Flatten(inp map[string]interface{}) map[string]interface{} {
	flatmap := make(map[string]interface{})
	flatten(inp, flatmap, "")
	return flatmap
}

func flatten(inp interface{}, flatmap map[string]interface{}, prefix string) {
	switch v := reflect.ValueOf(inp); v.Kind() {
	case reflect.Map:
		for key, val := range v.Interface().(map[string]interface{}) {
			flatten(val, flatmap, prefix+key+".")
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			flatten(v.Index(i).Interface(), flatmap, prefix+strconv.Itoa(i)+".")
		}
	default:
		flatmap[prefix[:len(prefix)-1]] = inp
	}
}
