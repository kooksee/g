package utils

import (
	"github.com/kooksee/g/assert"
	"net"
	"time"
)

func IpAddress() string {
	addrs, err := net.InterfaceAddrs()
	assert.Err(err, "net.InterfaceAddrs error")

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func If(b bool, tv, fv interface{}) interface{} {
	if b {
		return tv
	}
	return fv
}

func NowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TodayTimestamp() uint64 {
	n := time.Now().Unix()
	return uint64(n - n%(24*60*60) - 8*60*60 - 1)
}
