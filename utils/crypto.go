package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func GetFileSHA1(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	hash := sha1.New()

	if _, err := io.Copy(hash, f); err != nil {
		return ""
	}

	hib := hash.Sum(nil)[:20]

	return hex.EncodeToString(hib)
}

func ExistsAndValid(path, hash string) bool {
	if !FileExists(path) {
		return false
	}
	if GetFileSHA1(path) != hash {
		return false
	}
	return true
}
