package jun_util

import (
	"math/rand"
	"reflect"
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
