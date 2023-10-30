package helper

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

func MD5(p []byte) string {
	hash := md5.New()
	if _, err := hash.Write(p); err != nil {
		log.Fatalln(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
