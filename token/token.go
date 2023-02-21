/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : token.go
*   coder: zemanzeng
*   date : 2021-09-27 20:04:52
*   desc : 业务鉴权入口
*
================================================================*/

package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/wenruo95/jerry/aes"
	pb "github.com/wenruo95/jerry/token/proto"
	"google.golang.org/protobuf/proto"
)

const (
	// sdk版本号
	SDKVersion = "1.0.0"

	// sdk加密版本1
	EncodeVersion01 = 0x01

	// sdk支持的最大Encode版本
	SDKSupportMaxEncodeVersion = 0x01

	// 密钥默认过期时间 1天
	DefaultExpireTime = 24 * 3600
)

var (
	AuthErrSDKNeedUpdate = errors.New("sdk version too lower.")
)

// BizMessage 业务鉴权消息实现接口
type BizMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(buff []byte) error
	GetSecretKey(secretId uint32) string
}

// Encode 加密接口
// 请求参数:
//   - secretId  密钥ID
//   - secretKey 密钥
//   - message   加密数据
//   - expire    过期时间
//
// 返回参数:
//   - crypted	加密后的字段，鉴权字符串
//   - leftTime  剩余过期时间 单位秒
//   - error     错误信息
func Encode(secretId uint32, message BizMessage, expire int64) (crypted string, leftTime int64, err error) {
	tokenInfo := new(pb.BusinessToken)
	tokenInfo.Version = EncodeVersion01
	tokenInfo.SecretId = secretId
	if expire <= 0 {
		expire = DefaultExpireTime
	}
	tokenInfo.Ctime = time.Now().Unix()
	tokenInfo.Expire = tokenInfo.Ctime + expire

	switch tokenInfo.Version {
	case EncodeVersion01:
		messageBuff, err := message.Marshal()
		if err != nil {
			return "", 0, err
		}

		secretKey := message.GetSecretKey(secretId)
		if len(secretKey) == 0 {
			return "", 0, fmt.Errorf("secret[%v] key not exist", tokenInfo.SecretId)
		}

		data, err := aes.AesCBCEncrypt(messageBuff, []byte(secretKey))
		if err != nil {
			return "", 0, err
		}
		tokenInfo.Data = data

	default:
		return "", 0, fmt.Errorf("%w [sdk_version:%v encode_version:%v]",
			AuthErrSDKNeedUpdate, SDKVersion, tokenInfo.Version)
	}

	tokenBuff, err := proto.Marshal(tokenInfo)
	if err != nil {
		return "", 0, err
	}
	return base64.URLEncoding.EncodeToString(tokenBuff), expire, nil
}

// Decode 加密接口
// 请求参数:
//
//	-0. crypted 加密后的鉴权字符串
//	-1. message 业务信息 实现了序列化和反序列化接口(建议使用pb序列化) 数据反序列化到message中
//
// 返回参数:
//
//	-0. expire 剩余有效时间 单位:秒
//	-1. error 错误信息
func Decode(crypted string, message BizMessage) (leftTime int64, err error) {
	tokenBuff, err := base64.URLEncoding.DecodeString(crypted)
	if err != nil {
		return 0, err
	}

	tokenInfo := new(pb.BusinessToken)
	if err := proto.Unmarshal([]byte(tokenBuff), tokenInfo); err != nil {
		return 0, err
	}

	secretKey := message.GetSecretKey(tokenInfo.SecretId)
	if len(secretKey) == 0 {
		return 0, fmt.Errorf("secret[%v] key not exist", tokenInfo.SecretId)
	}

	switch tokenInfo.Version {
	case EncodeVersion01:
		messageBuff, err := aes.AesCBCDecrypt(tokenInfo.Data, []byte(secretKey))
		if err != nil {
			return 0, err
		}
		if err := message.Unmarshal(messageBuff); err != nil {
			return 0, err
		}

	default:
		return 0, fmt.Errorf("%w [sdk_version:%v encode_version:%v]",
			AuthErrSDKNeedUpdate, SDKVersion, tokenInfo.Version)

	}

	leftTime = tokenInfo.Expire - time.Now().Unix()
	if leftTime < 0 {
		return leftTime, fmt.Errorf("token has create:%v left:%v", tokenInfo.Ctime, leftTime)
	}
	return leftTime, nil
}
