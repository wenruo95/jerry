/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : misc.go
*   coder: zemanzeng
*   date : 2021-02-08 09:25:13
*   desc : 杂项
*
================================================================*/

package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
)

func ToJsonStr(v interface{}) string {
	if v == nil {
		return ""
	}

	buff, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(buff)
}

func GetBytes(v interface{}) []byte {
	buff := &bytes.Buffer{}
	binary.Write(buff, binary.BigEndian, v)
	return buff.Bytes()
}

func UUID() string {
	return uuid.New().String()
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func FileMD5(file *os.File) string {
	h := md5.New()
	io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}
