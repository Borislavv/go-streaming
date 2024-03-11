package ruid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func RequestUniqueID(r *http.Request) string {
	// machine hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		hostname = "unknown"
	}

	// process identifier
	pid := os.Getpid()

	// hash of commit of go version
	goVersion := runtime.Version()

	// timestamp with nanoseconds
	timestamp := time.Now().String()

	// random part of hash
	randomBytes := make([]byte, 16)
	if _, err = rand.Read(randomBytes); err != nil {
		log.Println(err)
	}
	randomString := hex.EncodeToString(randomBytes)

	strReqID := fmt.Sprintf(
		"%v%v%v%d%v%v%v%v%d%v",
		r.URL.String(),
		r.RemoteAddr,
		r.Method,
		r.ContentLength,
		r.Proto,
		timestamp,
		goVersion,
		hostname,
		pid,
		randomString,
	)

	hash := md5.New()
	if _, err = hash.Write([]byte(strReqID)); err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
