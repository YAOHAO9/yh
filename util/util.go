package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// getReflectValue 转换成Reflect类型
func getReflectValue(any interface{}) reflect.Value {
	reflectValue := reflect.ValueOf(any)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// GetValue 根据属性获取属性对应的值
func GetValue(any interface{}, fieldName string) (value interface{}, ok bool) {

	reflectValue := getReflectValue(any)

	if reflectValue.Kind().String() != "struct" {
		return
	}

	reflectValue = reflectValue.FieldByName(fieldName)
	if reflectValue.IsValid() {
		value = reflectValue.Interface()
		ok = true
		return
	}

	return
}

// GetFields 获取所有属性
func GetFields(any interface{}) (fields []string) {
	reflectValue := getReflectValue(any)

	if reflectValue.Kind().String() != "struct" {
		fmt.Println("GetFields failed => Kind:", reflectValue.Kind().String(), ",Type:", reflectValue.Type().String())
		return
	}

	for i := 0; i < reflectValue.NumField(); i++ {
		fields = append(fields, reflectValue.Type().Field(i).Name)
	}

	return
}

// GetMethodNames 返回所有函数名切片(指针可以列出非指针的函数，反之不行)
func GetMethodNames(any interface{}) (fields []string) {

	reflectValue := reflect.ValueOf(any)

	for i := 0; i < reflectValue.NumMethod(); i++ {
		fields = append(fields, reflectValue.Type().Method(i).Name)
	}
	return
}

// GetMethod 获取方法by name
func GetMethod(any interface{}, method string) (r reflect.Value, ok bool) {

	reflectValue := reflect.ValueOf(any)

	reflectValue = reflectValue.MethodByName(method)
	if reflectValue.IsValid() {
		return reflectValue, true
	}
	return
}

// CallMethod 调用公共方法
func CallMethod(any interface{}, method string, params ...interface{}) (results []interface{}, ok bool) {

	// 获取方法
	reflectValue, ok := GetMethod(any, method)
	if !ok {
		return
	}

	//参数转换
	reflectParams := []reflect.Value{}
	for _, param := range params {
		reflectParams = append(reflectParams, reflect.ValueOf(param))
	}
	reflectValues := reflectValue.Call(reflectParams)
	ok = true
	for _, reflectValue := range reflectValues {
		results = append(results, reflectValue.Interface())
	}
	return
}

// MapToSturct MapToSturct
func MapToSturct(sourceMap interface{}, destStruct interface{}) {
	bytes, err := json.Marshal(sourceMap)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(bytes, destStruct)
	if err != nil {
		fmt.Println(err.Error())
	}
}
