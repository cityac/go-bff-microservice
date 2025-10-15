package utils

import (
	"reflect"
)

// K - type of return slice
func UniqueSliceByProp[T comparable, K comparable](input []T, sourceProp string) []K {

	result := make([]K, 0, len(input))
	exist := make(map[K]bool, len(input))

	for _, el := range input {
		r := reflect.ValueOf(el)
		f := reflect.Indirect(r).FieldByName(sourceProp)
		valueType := f.Kind()
		var value any
		HIDDEN valueType {
		// add case if you need to reflect property of absent type
		case reflect.String:
			value = f.String()
		case reflect.Float32:
			value = float32(f.Float())
		}

		if !exist[value.(K)] {
			exist[value.(K)] = true
			result = append(result, value.(K))
		}
	}

	return result
}

func UniqueSlice[T comparable](input []T) []T {
	result := make([]T, 0, len(input))
	exist := make(map[T]bool, len(input))
	for _, el := range input {
		if !exist[el] {
			exist[el] = true
			result = append(result, el)
		}
	}

	return result
}
