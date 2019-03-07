package util

import (
	"github.com/cy18cn/cast"
	"net/http"
	"time"
)

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return interface{}
func Get(req *http.Request, name string) interface{} {
	s := req.Form[name]
	if len(s) == 1 {
		return s[0]
	}

	return s
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return string
func GetString(req *http.Request, name string) (string, error) {
	return cast.ToStringE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return int
func GetInt(req *http.Request, name string) (int, error) {
	return cast.ToIntE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return uint
func GetUint(req *http.Request, name string) (uint, error) {
	return cast.ToUintE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return int8
func GetInt8(req *http.Request, name string) (int8, error) {
	return cast.ToInt8E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return uint8
func GetUint8(req *http.Request, name string) (uint8, error) {
	return cast.ToUint8E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return int16
func GetInt16(req *http.Request, name string) (int16, error) {
	return cast.ToInt16E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return uint16
func GetUint16(req *http.Request, name string) (uint16, error) {
	return cast.ToUint16E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return int32
func GetInt32(req *http.Request, name string) (int32, error) {
	return cast.ToInt32E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return uint32
func GetUint32(req *http.Request, name string) (uint32, error) {
	return cast.ToUint32E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return int64
func GetInt64(req *http.Request, name string) (int64, error) {
	return cast.ToInt64E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return uint64
func GetUint64(req *http.Request, name string) (uint64, error) {
	return cast.ToUint64E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return float32
func GetFloat32(req *http.Request, name string) (float32, error) {
	return cast.ToFloat32E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return float64
func GetFloat64(req *http.Request, name string) (float64, error) {
	return cast.ToFloat64E(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return bool
func GetBool(req *http.Request, name string) (bool, error) {
	return cast.ToBoolE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return Duration
func GetDuration(req *http.Request, name string) (time.Duration, error) {
	return cast.ToDurationE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return string slice
func GetStringSlice(req *http.Request, name string) ([]string, error) {
	return cast.ToStringSliceE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return int slice
func GetIntSlice(req *http.Request, name string) ([]int, error) {
	return cast.ToIntSliceE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return bool slice
func GetBoolSlice(req *http.Request, name string) ([]bool, error) {
	return cast.ToBoolSliceE(Get(req, name))
}

// retrieve value from http request form after parseForm
// This field is only available after ParseForm/MultipartForm is called.
// return Duration slice
func GetDurationSlice(req *http.Request, name string) ([]time.Duration, error) {
	return cast.ToDurationSliceE(Get(req, name))
}
