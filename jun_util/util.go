package jun_util

import (
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() // 获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   // 获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		// 同类型
		valType := vTypeOfT.Field(i).Type

		if ok := bVal.FieldByName(name).IsValid() && bVal.FieldByName(name).Type() == valType; ok {
			bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
		}
	}
}
func Shuffle(vals []int32) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}
func InArray(arr []int32, target int32) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}
func GetHostIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip
}
func GetIpPorts(start, nums int) []int {
	// 根据接收参数个数，定义动态数组，
	var ip_ports []int
	for i := 0; i < nums; i++ {
		ip_ports = append(ip_ports, start+i)
	}
	return ip_ports
}

func CheckPorts(ip_ports []int) (port string) {
	for _, ip_port := range ip_ports {
		conn, err := net.DialTimeout("tcp", strconv.Itoa(ip_port), 3*time.Second)
		if err != nil {
			port = strconv.Itoa(ip_port)
			return
		}
		if conn != nil {
			conn.Close()
			continue
		} else {
			port = strconv.Itoa(ip_port)
			return
		}

	}
	return port
}
