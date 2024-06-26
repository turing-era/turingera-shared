package cutils

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/thinkeridea/go-extend/exnet"
	"google.golang.org/grpc/metadata"
)

// Str2Bytes string转[]byte无拷贝
func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2Str byte数组直接转成string对象，不发生内存copy
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Uint82Str uint8数组直接转成string对象,不发生内存copy
func Uint82Str(u []uint8) string {
	return *(*string)(unsafe.Pointer(&u))
}

// Obj2Json 将对象序列化成字符串
func Obj2Json(obj interface{}) string {
	bdata, err := jsoniter.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bdata)
}

// RandString 随机字符串
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

// GetLocalIP 获取本机ip
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
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

// GetClientIp 获取ctx中header client ip
func GetClientIp(ctx context.Context) string {
	ip := metadata.ValueFromIncomingContext(ctx, "x-forwarded-for")[0]
	for _, ip = range strings.Split(ip, ",") {
		if ip = strings.TrimSpace(ip); ip != "" && !exnet.HasLocalIPAddr(ip) {
			return ip
		}
	}
	ip = metadata.ValueFromIncomingContext(ctx, "x-real-ip")[0]
	if ip = strings.TrimSpace(ip); ip != "" && !exnet.HasLocalIPAddr(ip) {
		return ip
	}
	return ""
}

// GetTodayZeroTime 获取当日0点时间戳
func GetTodayZeroTime() int32 {
	t := time.Now()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int32(newTime.Unix())
}
