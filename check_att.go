package helpers

import (
	"reflect"
	"strings"
)

// CheckInTag represent the helpers for check attribute inside of struct
func CheckInTag(v interface{}, filterField, tag string) []string {
	filterField = strings.Replace(filterField, " ", "", -1)
	field := strings.Split(filterField, ",")

	out := []string{}
	mapField := make(map[string]bool)
	for _, f := range field {
		mapField[f] = true
	}
	val := reflect.ValueOf(v)
	for i := 0; i < val.Type().NumField(); i++ {
		if tag, ok := val.Type().Field(i).Tag.Lookup(tag); ok {
			if _, ok := mapField[tag]; ok {
				if tag != "-" {
					out = append(out, tag)
				}
			}
		}
	}

	return out
}

// CheckMatchTag represent the helpers for check attribute match of struct
func CheckMatchTag(v interface{}, field, tag string) string {
	var out string
	val := reflect.ValueOf(v)
	for i := 0; i < val.Type().NumField(); i++ {
		if tag, ok := val.Type().Field(i).Tag.Lookup(tag); ok {
			if tag == field {
				if tag != "-" {
					out = tag
				}
			}
		}
	}

	return out
}
