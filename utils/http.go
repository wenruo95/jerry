/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : http.go
*   coder: zemanzeng
*   date : 2021-09-27 19:24:41
*   desc : http request
*
================================================================*/

package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func SendHttpRequest(httpURL string, request interface{}, response interface{}) error {
	src, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, httpURL, bytes.NewBuffer(src))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dst, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(dst, response); err != nil {
		return err
	}

	return nil
}
