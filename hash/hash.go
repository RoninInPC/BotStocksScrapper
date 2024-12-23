package hash

import (
	"crypto/md5"
	"fmt"
)

type Hash func(string) string

func MD5(str string) string {
	data := []byte(str)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	return hash
}
