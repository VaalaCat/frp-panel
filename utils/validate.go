package utils

const (
	whiteListChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
)

func IsClientIDPermited(clientID string) bool {
	if len(clientID) == 0 {
		return false
	}

	chrMap := make(map[rune]bool)
	for _, chr := range whiteListChar {
		chrMap[chr] = true
	}

	for _, chr := range clientID {
		if !chrMap[chr] {
			return false
		}
	}

	return true
}

func MakeClientIDPermited(clientID string) string {
	input := []rune(clientID)
	output := input
	chrMap := make(map[rune]bool)
	for _, chr := range whiteListChar {
		chrMap[chr] = true
	}
	for idx, chr := range input {
		if !chrMap[chr] {
			output[idx] = '-'
		}
	}
	return string(output)
}
