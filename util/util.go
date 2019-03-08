package util

import (
	"bytes"
	"reflect"
)

func IsStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func BytesCombine(b ...[]byte) []byte {
	return bytes.Join(b, []byte(""))
}
