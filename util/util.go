package util

import "reflect"

func IsStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}
