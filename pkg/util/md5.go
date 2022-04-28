package util

import (
	"crypto/md5" // nolint: gosec
)

// Md5 returns the MD5 checksum of the given data.
func Md5(data []byte) []byte {
	hash := md5.New() // nolint: gosec
	hash.Write(data)
	return hash.Sum(nil)
}

// Md5String returns the MD5 checksum of the given data as a hex string.
func Md5String(str string) string {
	return string(Md5([]byte(str)))
}
