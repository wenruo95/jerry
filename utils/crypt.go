/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : crypt.go
*   coder: zemanzeng
*   date : 2021-02-08 09:43:50
*   desc : 加密通用封装
*
================================================================*/

package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func MD5Data(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5File(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	h := md5.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Sha1Data(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
