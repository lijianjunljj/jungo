package jun_util

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
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

func Shuffle[T any](vals []T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}
func InArray(arr interface{}, target interface{}) bool {
	switch v := arr.(type) {
	case []int:
		for _, item := range v {
			if item == target {
				return true
			}
		}
	case []string:
		for _, item := range v {
			if item == target {
				return true
			}
		}
	// 可以根据需要添加更多类型的处理
	default:
		fmt.Println("不支持的数组类型")
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

func GetIPV4() string {
	resp, err := http.Get("https://ipv4.netarm.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func GetIPV6() string {
	resp, err := http.Get("https://ipv6.netarm.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
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
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(ip_port), 3*time.Second)

		//fmt.Println("err:", err)
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
