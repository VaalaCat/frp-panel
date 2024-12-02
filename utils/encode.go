package utils

import (
	"encoding/base64"
)

func EncodeBase64(data string) string {
	encodedStr := base64.StdEncoding.EncodeToString([]byte(data))
	return encodedStr
}
