/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : aes.go
*   coder: zemanzeng
*   date : 2021-05-25 16:18:50
*   desc : aes加密解密
*
================================================================*/

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
)

/*
加密模式: `CBC`
填充: `PKCS7Padding`
密钥: 取`AESKey`前16位
偏移量: 和密钥相同，取`AESKey`前16位
密文编码: `HEX`
线上测试工具: <https://oktools.net/aes>
*/

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesCBCEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesCBCDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

const (
	AesKeyL16   = 16
	AesKeyL32   = 32
	AesKeyLRand = 0
)

func GenarateAesKey(n int) string {
	if n != AesKeyL16 && n != AesKeyL32 {
		switch rand.Int() % 2 {
		case 0:
			n = AesKeyL16
		case 1:
			n = AesKeyL32
		default:
		}
	}

	var key []rune
	var runes []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := 0; i < n; i++ {
		key = append(key, runes[rand.Int()%len(runes)])
	}
	return string(key)
}
