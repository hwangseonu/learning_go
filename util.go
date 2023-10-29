package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
)

func logHandler(_ http.ResponseWriter, req *http.Request) {
	log.Printf("[%s] %s", req.Method, req.URL.Path)
}

func hashString(data []byte) string {
	hash := sha256.New()
	hash.Write(data)

	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

//func compareHash(a, b []byte) bool {
//	return bytes.Equal(a, b)
//}
