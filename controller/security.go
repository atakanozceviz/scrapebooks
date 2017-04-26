package controller

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"
)

func Secret(s string) bool {
	check := make(map[string]string)
	t := time.Now().UTC().Unix()

	for i := 0; i <= 5; i++ {
		secret := []byte("aratoz" + strconv.FormatInt(t-int64(i), 10))
		res := sha1.Sum(secret)
		check[hex.EncodeToString(res[:])] = ""

	}

	if _, ok := check[s]; ok {
		return true
	} else {
		return false
	}
}
