package utils

import (
	"context"
	"reflect"
	"strconv"
)

func GetValue[T any](c context.Context, key string) (T, bool) {
	val, ok := getValue[T](c, key)
	if !ok {
		return *new(T), false
	}
	v, ok := val.(T)
	if !ok {
		return *new(T), false
	}
	return v, true
}

func getValue[T any](c context.Context, key string) (any, bool) {
	val := c.Value(key)
	if val == nil {
		return *new(T), false
	}

	if reflect.TypeOf(*new(T)).Kind() == reflect.Int && reflect.TypeOf(val).Kind() == reflect.String {
		strconvInt, err := strconv.Atoi(val.(string))
		if err != nil {
			return *new(T), false
		}
		return strconvInt, true
	}

	v, ok := val.(T)
	if !ok {
		return *new(T), false
	}

	return v, true
}

func GetIntValueFromStr(c context.Context, key string) (int, bool) {
	val := c.Value(key)
	if val == nil {
		return 0, false
	}

	v, ok := val.(string)
	if !ok {
		return 0, false
	}

	strconvInt, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}

	return strconvInt, true
}
