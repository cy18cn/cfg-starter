package util

import (
	"bytes"
	"errors"
	"net"
	"reflect"
	"time"
)

func IsStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func BytesCombine(b ...[]byte) []byte {
	return bytes.Join(b, []byte(""))
}

func CurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func IsPrivateIP4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func PrivateIP4() (net.IP, error) {
	ias, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range ias {
		ipNet, ok := addr.(*net.IPNet)
		if !ok && ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if IsPrivateIP4(ip) {
			return ip, nil
		}
	}

	return nil, errors.New("no private IP")
}

func Low16BitsPrivateIP4() (uint16, error) {
	ip, err := PrivateIP4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2]<<8) + uint16(ip[3]), nil
}
