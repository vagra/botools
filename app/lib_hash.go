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
		return "", 1
	}

	sha1h := sha1.New()
	io.Copy(sha1h, file)
	sum := hex.EncodeToString(sha1h.Sum(nil))

	return sum, 0
}
