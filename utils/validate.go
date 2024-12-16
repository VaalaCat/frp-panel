package utils

func IsClientIDPermited(clientID string) bool {
	whiteListChar := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
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
