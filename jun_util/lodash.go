package jun_util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/olebedev/config"
	uuid "github.com/satori/go.uuid"
)

// Md5 将字符串进行MD5加密
func Md5(str string) string {
	b := []byte(str)
	s := fmt.Sprintf("%x", md5.Sum(b))
	return s
}

// UUID 生成UUID字符串
func UUID() string {
	s := uuid.Must(uuid.NewV4(), nil).String()
	// s = strings.Replace(s, "-", "", -1)
	return s
}

// Password 两次MD5加密生成密码
func Password(str string) string {
	if len(str) == 0 {
		return ""
	}
	return Md5(Md5(str))
}

// Capitalize 首字母大写
func Capitalize(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArr := []rune(str)
	if strArr[0] >= 97 && strArr[0] <= 122 {
		strArr[0] -= 32
	}
	return string(strArr)
}

// FirstToLower 首字母小写
func FirstToLower(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArr := []rune(str)
	if strArr[0] <= 97 || strArr[0] >= 122 {
		strArr[0] += 32
	}
	return string(strArr)
}

// FirstLetter 首字母
func FirstLetter(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArr := []rune(str)
	if strArr[0] <= 97 || strArr[0] >= 122 {
		strArr[0] += 32
	}
	return string(strArr[0])
}

// Bind 绑定入参到结构体或者MAP,object为指针
func Bind(ctx *gin.Context, object interface{}) (interface{}, error) {
	bindType := binding.Default(ctx.Request.Method, ctx.ContentType())
	if err := ctx.ShouldBindWith(object, bindType); err != nil {
		return nil, err
	}
	return object, nil
}

// Invoke 传入函数和参数调用函数返回函数执行结果
func Invoke(f interface{}, params ...interface{}) []reflect.Value {
	fv := reflect.ValueOf(f)
	realParams := make([]reflect.Value, len(params))
	for i, item := range params {
		realParams[i] = reflect.ValueOf(item)
	}
	rs := fv.Call(realParams)
	return rs
}

// Getenv 获取ENV环境变量
func Getenv() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	return env
}

// IsEmpty 判断变量是否为空
func IsEmpty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// Base64Encode 字符串转base64编码
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func UrlEncode(str string) string {
	return url.QueryEscape(str)
}
func UrlDecode(str string) string {
	res, err := url.QueryUnescape(str)
	if err != nil {
		return ""
	}
	return res
}

// Base64Decode base64编码转字符串
func Base64Decode(s string) (string, error) {
	ds, err := base64.StdEncoding.DecodeString(s)
	return string(ds), err
}

// JSONEncode MAP/结构体转JSON字符串
func JSONEncode(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JSONDecode JSON字符串转结构体/MAP
func JSONDecode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

// FloatFixed 小数点后 n 位四舍五入
func FloatFixed(val float64, n int) float64 {
	shift := math.Pow(10, float64(n))
	fv := 0.0000000001 + val //对浮点数产生.xxx999999999 计算不准进行处理
	return math.Floor(fv*shift+.5) / shift
}

// FloatRound 小数点后 n 位舍去
func FloatRound(val float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n+1)+"f", val)
	temp := strings.Split(floatStr, ".")
	var newFloat string
	if len(temp) < 2 || n >= len(temp[1]) {
		newFloat = floatStr
	} else {
		newFloat = temp[0] + "." + temp[1][:n]
	}
	inst, _ := strconv.ParseFloat(newFloat, 64)
	return inst
}

// ParseYAML 从YAML配置文件中获取配置对象
func ParseYAML(filepath string) (*config.Config, error) {
	cfg, err := config.ParseYamlFile(filepath)
	return cfg, err
}

// ParseJSON 从JSON配置文件中获取配置对象
func ParseJSON(filepath string) (*config.Config, error) {
	cfg, err := config.ParseJsonFile(filepath)
	return cfg, err
}

// GetIP 获取本地IP地址
func GetIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := GetIPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("connected to the network?")
}

// GetIPFromAddr 从地址信息中获取IP
func GetIPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}

	return ip
}

// JoinURL 连接网址
func JoinURL(urls ...string) string {
	url := ""
	for _, val := range urls {
		url = path.Join(url, val)
	}
	return strings.Replace(url, ":/", "://", 1)
}

// GenCode 生成N位随机数字字符串
func GenCode(width int) string {
	rand.Seed(time.Now().UnixNano())
	var result strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&result, "%d", rand.Intn(10))
	}
	return result.String()
}

// GenStr 生成N位随机字符串
func GenStr(width int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, width)
	for i := 0; i < width; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(PanicTrace(e))
		}
	}()
	f()
	return
}

// PanicTrace panic调用链跟踪
func PanicTrace(err interface{}) string {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)

	return fmt.Sprintf("panic: %v %s", err, stackBuf[:n])
}

// InStringSlice 判断字符串是否在切片中
func InStringSlice(slice []string, element string) bool {
	element = strings.TrimSpace(element)
	for _, v := range slice {
		if strings.TrimSpace(v) == element {
			return true
		}
	}

	return false
}

// CurrentPath 当前路径
func CurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// RetryTimes 重试
func RetryTimes(tryTimes int, sleep time.Duration, callback func() error) (err error) {
	for i := 1; i <= tryTimes; i++ {
		err = callback()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)
	}
	return err
}

// RandPicker 随机取N个值
func RandPicker(origin []string, count int) []string {
	tmpOrigin := make([]string, len(origin))
	copy(tmpOrigin, origin)
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(tmpOrigin), func(i int, j int) {
		tmpOrigin[i], tmpOrigin[j] = tmpOrigin[j], tmpOrigin[i]
	})

	result := make([]string, 0, count)
	for index, value := range tmpOrigin {
		if index == count {
			break
		}
		result = append(result, value)
	}
	return result
}

// Decimal 64位浮点数保留2位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// Sum 求和
func Sum(results []float64) float64 {
	var sum float64
	for _, v := range results {
		sum += v
	}
	return Decimal(sum)
}
