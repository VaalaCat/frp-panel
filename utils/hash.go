package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func MD5[T []byte | string](input T) string {
	data := []byte(input)
	hash := md5.Sum(data)
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func SHA1(input string) string {
	hash := sha1.Sum([]byte(input))
	return fmt.Sprintf("%x", hash)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MatchSecret compares fixed-length digests so token length is not a timing signal.
func MatchSecret(provided, expected string) bool {
	a, b := sha256.Sum256([]byte(provided)), sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1 && provided != "" && expected != ""
}

// Use the normal cost without generating a bcrypt hash during package initialization.
const dummyPasswordHash = "$2a$14$iOcKdM4wLfLPzHOcGPosduTT/9qHvjnS8jAZI/4INNhdIZcPFhk2u"

// CheckDummyPassword performs the normal bcrypt work for a nonexistent account.
func CheckDummyPassword(password string) {
	CheckPasswordHash(password, dummyPasswordHash)
}
