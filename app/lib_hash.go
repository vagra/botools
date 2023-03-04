package app

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func GetSHA1(path string) (string, int8) {

	file, err := os.Open(path)
	if err != nil {
		return "", 10
	}

	sha1h := sha1.New()
	_, err = io.Copy(sha1h, file)
	if err != nil {
		return "", 20
	}

	sum := hex.EncodeToString(sha1h.Sum(nil))

	return sum, 0
}
