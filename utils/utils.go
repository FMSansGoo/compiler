package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func LogInfo(msg string, args ...interface{}) {
	str := fmt.Sprintf("[INFO] %s", msg)
	if len(args) > 0 {
		params := make([]string, 0, len(args))
		for _, arg := range args {
			value := reflect.ValueOf(arg)
			name := reflect.TypeOf(arg).Name()
			params = append(params, fmt.Sprintf("%s=%+v", name, value))
		}
		str = fmt.Sprintf("[INFO] %s: %s", msg, strings.Join(params, ","))
	}
	fmt.Println(str)
}

func LogError(msg string, args ...interface{}) {
	str := fmt.Sprintf("[ERROR] %s", msg)
	if len(args) > 0 {
		params := make([]string, 0, len(args))
		for _, arg := range args {
			value := reflect.ValueOf(arg)
			name := reflect.TypeOf(arg).Name()
			params = append(params, fmt.Sprintf("%s=%v", name, value))
		}
		str = fmt.Sprintf("[ERROR] %s: %s", msg, strings.Join(params, ","))
	}
	fmt.Println(str)
	panic(str)
}

func LogErrorFormat(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] %s: %s\n", msg, fmt.Sprint(args...))
}

func InStringSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
