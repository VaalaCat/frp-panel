package utils

import "encoding/json"

func MarshalForJson(v any) string {
	ret, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(ret)
}
