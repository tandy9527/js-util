package tools

import (
	"fmt"
	"reflect"
)

// GetNested 获取嵌套的值,类似树结构数据
// m: map[string]any
// path: []string 路径
// T: any 返回的值类型

// GetNested 根据路径获取配置，并自动做类型转换
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
			panic(fmt.Sprintf(
				"[GetNested] config key not found: '%v'\nFull path: %v\nMap keys at this level: %v",
				p, path[:i+1], reflect.ValueOf(currMap).MapKeys(),
			))
		}

		curr = val
	}

	// 自动转换类型
	curr = autoConvert(curr)

	// 类型检查
	t, ok := curr.(T)
	if !ok {
		panic(fmt.Sprintf(
			"[GetNested] config type error:\n  Expected: %T\n  Got: %T (value=%v)\n  Path: %v",
			*new(T), curr, curr, path,
		))
	}
	return t
}

// autoConvert 自动转换 YAML 解析后的类型
func autoConvert(v any) any {
	switch val := v.(type) {

	case float64: // YAML 默认数字是 float64
		return int(val)

	case []any:
		if len(val) == 0 {
			return []int{}
		}
		// 判断是二维数组还是一维数组
		switch val[0].(type) {
		case []any: // [][]int
			result := make([][]int, 0, len(val))
			for _, item := range val {
				subArr, ok := item.([]any)
				if !ok {
					panic(fmt.Sprintf("[autoConvert] expected []any in [][]int but got %T", item))
				}
				row := make([]int, 0, len(subArr))
				for _, n := range subArr {
					row = append(row, toInt(n))
				}
				result = append(result, row)
			}
			return result
		default: // []int
			result := make([]int, 0, len(val))
			for _, n := range val {
				result = append(result, toInt(n))
			}
			return result
		}

	case map[interface{}]interface{}: // key 可能是数字
		strMap := make(map[string]any)
		for k, v2 := range val {
			strKey := fmt.Sprintf("%v", k)
			strMap[strKey] = autoConvert(v2)
		}
		return strMap

	case map[string]any:
		for k, v2 := range val {
			val[k] = autoConvert(v2)
		}
		return val

	default:
		return val
	}
}

// toInt 将任何数字类型转换为 int
func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		panic(fmt.Sprintf("[toInt] cannot convert type %T to int (value=%v)", v, v))
	}
}
