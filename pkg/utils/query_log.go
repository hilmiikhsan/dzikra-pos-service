package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func flattenArgs(args []interface{}) []interface{} {
	var flat []interface{}
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				flat = append(flat, val.Index(i).Interface())
			}
		} else {
			flat = append(flat, arg)
		}
	}
	return flat
}

func QueryLog(query string, args ...interface{}) {
	flatArgs := flattenArgs(args)

	if strings.Contains(query, "$1") {
		for i, v := range flatArgs {
			placeholder := fmt.Sprintf("$%d", i+1)
			query = strings.Replace(query, placeholder, fmt.Sprintf("'%v'", v), 1)
		}
	} else if strings.Contains(query, "?") {
		for _, v := range flatArgs {
			query = strings.Replace(query, "?", fmt.Sprintf("'%v'", v), 1)
		}
	}
	fmt.Println(query)
}
