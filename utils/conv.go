package utils

import (
	"strconv"
)

func Str2Int64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func Str2Int64Default(str string, intVal int64) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return intVal
	}
	return num
}

func ToStr(any interface{}) string {
	if any == nil {
		return ""
	}

	if str, ok := any.(string); ok {
		return str
	}

	return ""
}

func IsInteger(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
