package tools

import (
	"fmt"
	"reflect"
)

// GetNested 获取嵌套的值,类似树结构数据
// m: map[string]any
// path: []string 路径
// T: any 返回的值类型
// 支持任意嵌套结构
func GetNested[T any](m map[string]any, path ...string) T {
	if len(path) == 0 {
		panic("[GetNested] empty path: at least one key is required")
	}

	var curr any = m
	for i, p := range path {
		currMap, ok := curr.(map[string]any)
		if !ok {
			panic(fmt.Sprintf(
				"[GetNested] invalid config path at '%v': expected map[string]any but got %T\nPath so far: %v",
				p, curr, path[:i+1],
			))
		}

		val, exists := currMap[p]
		if !exists {
			// 尝试数字 key 转 string
			for k, v := range currMap {
				if fmt.Sprintf("%v", k) == p {
					val = v
					exists = true
					break
				}
			}
		}

		if !exists {
			panic(fmt.Sprintf(
				"[GetNested] config key not found: '%v'\nFull path: %v\nMap keys: %v",
				p, path[:i+1], reflect.ValueOf(currMap).MapKeys(),
			))
		}

		curr = val
	}

	// 根据 T 类型递归转换
	tType := reflect.TypeOf((*T)(nil)).Elem()
	result := convertToType(curr, tType)
	return result.(T)
}

// convertToType 通用递归转换函数
func convertToType(v any, t reflect.Type) any {
	switch t.Kind() {
	case reflect.Int:
		return toInt(v)
	case reflect.Float64:
		return toFloat64(v)
	case reflect.String:
		return fmt.Sprintf("%v", v)
	case reflect.Bool:
		b, ok := v.(bool)
		if !ok {
			panic(fmt.Sprintf("[convertToType] expected bool but got %T", v))
		}
		return b

	case reflect.Slice:
		arr, ok := v.([]any)
		if !ok {
			panic(fmt.Sprintf("[convertToType] expected []any but got %T", v))
		}
		res := reflect.MakeSlice(t, len(arr), len(arr))
		for i, elem := range arr {
			res.Index(i).Set(reflect.ValueOf(convertToType(elem, t.Elem())))
		}
		return res.Interface()

	case reflect.Map:
		m, ok := v.(map[string]any)
		if !ok {
			panic(fmt.Sprintf("[convertToType] expected map[string]any but got %T", v))
		}
		res := reflect.MakeMap(t)
		for k, val := range m {
			keyVal := reflect.ValueOf(k)
			valVal := reflect.ValueOf(convertToType(val, t.Elem()))
			res.SetMapIndex(keyVal, valVal)
		}
		return res.Interface()

	default:
		panic(fmt.Sprintf("[convertToType] unsupported kind: %v", t.Kind()))
	}
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		panic(fmt.Sprintf("[toInt] cannot convert %T to int", v))
	}
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		panic(fmt.Sprintf("[toFloat64] cannot convert %T to float64", v))
	}
}
