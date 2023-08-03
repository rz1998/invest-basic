package invest

import (
	"fmt"
	"reflect"
)

// SumStructs 对结构体序列的指定字段求和
func SumStructs(nameFieldSums []string, structs []interface{}) map[string]float64 {
	result := make(map[string]float64)
	var valData reflect.Value
	for _, s := range structs {
		typeData := reflect.TypeOf(s)
		if typeData.Kind() == reflect.Ptr {
			valData = reflect.ValueOf(s).Elem()
		} else {
			valData = reflect.ValueOf(s)
		}
		for _, name := range nameFieldSums {
			field := valData.FieldByName(name)
			var val float64
			switch field.Kind() {
			case reflect.Float32, reflect.Float64:
				val = field.Float()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val = float64(field.Int())
			default:
				fmt.Printf("unhandled type %v\n", field.Kind())
			}
			result[name] += val
		}
	}
	return result
}

// WAvgStruct 加权平均序列
func WAvgStruct(nameWeight, nameData string, dataIn []interface{}) float64 {
	// 求和
	mapSum := SumStructs([]string{nameWeight}, dataIn)
	weightTotal := mapSum[nameWeight]

	var result float64
	var valData reflect.Value
	for _, s := range dataIn {
		typeData := reflect.TypeOf(s)
		if typeData.Kind() == reflect.Ptr {
			valData = reflect.ValueOf(s).Elem()
		} else {
			valData = reflect.ValueOf(s)
		}
		// 数据
		fieldD := valData.FieldByName(nameData)
		var valD float64
		switch fieldD.Kind() {
		case reflect.Float32, reflect.Float64:
			valD = fieldD.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valD = float64(fieldD.Int())
		default:
			fmt.Printf("unhandled type %v\n", fieldD.Kind())
		}
		// 权重
		fieldW := valData.FieldByName(nameWeight)
		var valW float64
		switch fieldW.Kind() {
		case reflect.Float32, reflect.Float64:
			valW = fieldW.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valW = float64(fieldW.Int())
		default:
			fmt.Printf("unhandled type %v\n", fieldW.Kind())
		}
		// 计算加权
		result += valD * valW / weightTotal
	}
	return result
}
