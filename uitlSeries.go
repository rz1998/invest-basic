package invest

import (
	"container/list"
	"fmt"
	"reflect"
	"time"
)

// SliceTimeSeries 分拆时间序列数据
func SliceTimeSeries(nameTimestamp string, duration time.Duration, timeSeries []interface{}) [][]interface{} {
	var result [][]interface{}
	interval := duration.Milliseconds()
	var timeFloor, timeCeil int64
	dataList := list.New()
	var valData reflect.Value
	for _, data := range timeSeries {
		typeData := reflect.TypeOf(data)
		if typeData.Kind() == reflect.Ptr {
			valData = reflect.ValueOf(data).Elem()
		} else {
			valData = reflect.ValueOf(data)
		}
		timestamp := valData.FieldByName(nameTimestamp).Int()
		if timestamp >= timeCeil {
			if dataList != nil && dataList.Len() > 0 {
				// 已有数据则转换为切片保存到结果里
				vals := make([]interface{}, dataList.Len())
				result = append(result, vals)
				index := 0
				for i := dataList.Front(); i != nil; i = i.Next() {
					vals[index] = i.Value
					index++
				}
				dataList.Init()
			}
			timeFloor = timestamp / interval * interval
			timeCeil = (timestamp/interval + 1) * interval
		}
		valData.FieldByName(nameTimestamp).SetInt(timeFloor)
		dataList.PushBack(data)
	}
	return result
}

// MergeTimeSeries 合并时间序列数据
func MergeTimeSeries(nameFieldSums []string, timeSeries []interface{}) interface{} {
	// 最后一条作为基础
	data := timeSeries[len(timeSeries)-1]
	// 只汇总需要汇总的字段
	sumMap := SumStructs(nameFieldSums, timeSeries)
	// 求和项放回最后一条
	typeData := reflect.TypeOf(data)
	var valData reflect.Value
	if typeData.Kind() == reflect.Ptr {
		valData = reflect.ValueOf(data).Elem()
	} else {
		valData = reflect.ValueOf(data)
	}
	for _, name := range nameFieldSums {
		field := valData.FieldByName(name)
		switch field.Kind() {
		case reflect.Float32, reflect.Float64:
			field.SetFloat(sumMap[name])
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(int64(sumMap[name]))
		default:
			fmt.Printf("unhandled type %v\n", field.Kind())
		}
	}
	return data
}
