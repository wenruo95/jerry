/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : http.go
*   coder: zemanzeng
*   date : 2022-01-17 11:59:11
*   desc : http常用封装
*
================================================================*/

package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadBody(r *http.Request, body interface{}) error {
	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(buff, body); err != nil {
		return err
	}
	return nil
}

func WriteBody(w http.ResponseWriter, body interface{}) error {
	buff, err := json.Marshal(body)
	if err != nil {
		if _, err := w.Write([]byte("internal error")); err != nil {
			return err
		}
		return err
	}

	if _, err := w.Write(buff); err != nil {
		return err
	}
	return nil
}
