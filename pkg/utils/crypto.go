package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// MD5 generates the MD5 hash of the input string
func MD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	signature := hex.EncodeToString(hash.Sum(nil))
	return signature
}

// SHA256 generates the SHA-256 hash of the input string
func SHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	signature := hex.EncodeToString(hash.Sum(nil))
	return signature
}

// HMAC256 generates the SHA-256 hash of the input string
func HMAC256(input, secret string) string {
	key := []byte(secret)
	msg := []byte(input)
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	signature := mac.Sum(nil)
	return hex.EncodeToString(signature)
}
