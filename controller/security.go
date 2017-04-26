package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
)

//Compares the key with FormValue("key")
func CheckKey(key string) bool {
	t := time.Now().UTC().Unix()
	check := make(map[string]string)
	//Makes client have 5 seconds to use key
	for i := 0; i <= 5; i++ {
		secret := []byte("aratoz" + strconv.FormatInt(t-int64(i), 10))
		res := sha1.Sum(secret)
		check[hex.EncodeToString(res[:])] = ""
	}
	//Returns the result
	if _, ok := check[key]; ok {
		return true
	} else {
		return false
	}
}
