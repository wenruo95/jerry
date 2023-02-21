/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : values.go
*   coder: zemanzeng
*   date : 2021-09-29 17:00:07
*   desc : values配置解析封装
*
================================================================*/

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Values 配置解析封装
type Values map[string]interface{}

func (values Values) GetValue(key string) (interface{}, error) {
	if values == nil {
		return nil, errors.New("values is nil")
	}
	if len(key) == 0 {
		return nil, errors.New("key length is zero")
	}

	iValue, exist := values[key]
	if !exist {
		return nil, errors.New("key:" + key + " not exist")
	}
	return iValue, nil
}

func (values Values) GetJson(key string, data interface{}) error {
	iValue, err := values.GetValue(key)
	if err != nil {
		return err
	}

	switch value := iValue.(type) {
	case string:
		return json.Unmarshal([]byte(value), data)
	case []byte:
		return json.Unmarshal(value, data)
	default:
		return fmt.Errorf("invalid type key:%v type:%T value:%v", key, value, value)
	}
}

func (values Values) GetYaml(key string, data interface{}) error {
	iValue, err := values.GetValue(key)
	if err != nil {
		return err
	}

	switch value := iValue.(type) {
	case string:
		return yaml.Unmarshal([]byte(value), data)
	case []byte:
		return yaml.Unmarshal(value, data)
	default:
		return fmt.Errorf("invalid type key:%v type:%T value:%v", key, value, value)
	}
}

func (values Values) GetList(key string) ([]string, error) {
	iValue, err := values.GetValue(key)
	if err != nil {
		return nil, err
	}

	switch value := iValue.(type) {
	case string:
		return strings.Fields(value), nil
	case []byte:
		return strings.Fields(string(value)), nil
	default:
		return nil, fmt.Errorf("invalid type key:%v type:%T value:%v", key, value, value)
	}
}

func (values Values) GetLists() (map[string][]string, error) {
	if values == nil {
		return nil, errors.New("values is nil")
	}

	datas := make(map[string][]string)
	for key, iValue := range values {
		switch value := iValue.(type) {
		case string:
			datas[key] = strings.Fields(value)
		case []byte:
			datas[key] = strings.Fields(string(value))
		default:
			return nil, fmt.Errorf("invalid type key:%v type:%T value:%v", key, value, value)
		}
	}
	return datas, nil
}
