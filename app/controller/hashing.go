package controller

import (
	"crypto/md5"
	"fmt"
)

func Hash(str string) (string, error) {
	data := []byte(str)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	return hash, nil
}
