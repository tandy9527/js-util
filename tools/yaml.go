package tools

import (
	"fmt"
	"os"
	"reflect"

	"github.com/tandy9527/js-util/logger"
	"gopkg.in/yaml.v3"
)

// Loadyaml 加载.yaml
func Loadyaml[T any](filePath string) *T {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic("failed to read config file: " + filePath)
	}
	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic("failed to unmarshal config file: " + filePath)
	}

	logger.Infof("[Loadyaml] load successful: %s", filePath)
	return &config
}

// GetNested 获取嵌套的值,类似树结构数据
// m: map[string]any
// path: []string 路径
// T: any 返回的值类型
func GetNested[T any](m map[string]any, path ...string) T {
	if len(path) == 0 {
		panic("[getNested] empty path: at least one key is required")
	}
	var curr any = m
	for i, p := range path {
		currMap, ok := curr.(map[string]any)
		if !ok {
			panic(fmt.Sprintf(
				"[getNested] invalid config path at '%v': expected map[string]any but got %T\n"+
					"Path so far: %v",
				p, curr, path[:i+1],
			))
		}

		val, exists := currMap[p]
		if !exists {
			panic(fmt.Sprintf(
				"[getNested] config key not found: '%v'\n"+
					"Full path: %v\n"+
					"Map keys at this level: %v",
				p, path[:i+1], reflect.ValueOf(currMap).MapKeys(),
			))
		}

		curr = val
	}

	t, ok := curr.(T)
	if !ok {
		panic(fmt.Sprintf(
			"[getNested] config type error:\n"+
				"  Expected: %T\n"+
				"  Got: %T (value=%v)\n"+
				"  Path: %v",
			*new(T), curr, curr, path,
		))
	}

	return t
}
