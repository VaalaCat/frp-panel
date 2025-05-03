package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUIDWithoutSeperator() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func GenerateUUID() string {
	return uuid.New().String()
}
