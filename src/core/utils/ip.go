package utils

import (
	"math"
	"net"
)

func IsIP(ip string) bool {
	address := net.ParseIP(ip)
	if nil == address {
		return false
	}

	return true
}

func IsPort(port int64) bool {
	if (port < 0) || (port > 65535) {
		return false
	}

	return true
}

func IsLegalTimeout(timeout int64) bool {
	if timeout <= 0 || timeout > math.MaxInt32 {
		return false
	}

	return true
}

func IsLegalNumber(n int) bool {
	if n < 0 || n > math.MaxInt32 {
		return false
	}

	return true
}

func IsLegalInterval(interval int64) bool {
	if interval < 0 || interval > math.MaxInt32 {
		return false
	}

	return true
}

func IsLegalStatisticsInterVal(interval int64) bool {
	if interval <= 0 || interval > math.MaxInt32 {
		return false
	}

	return true
}

