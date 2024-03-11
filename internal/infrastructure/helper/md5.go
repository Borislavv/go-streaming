package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(p []byte) string {
	hash := md5.New()
	if _, err := hash.Write(p); err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
