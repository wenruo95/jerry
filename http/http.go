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

	"google.golang.org/protobuf/proto"
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

func JsonBodyToPb(r *http.Request, pbData proto.Message) error {
	if params := r.URL.Query(); len(params) > 0 {
		paramsBuff, err := json.Marshal(params)
		if err != nil {
			return err
		}
		if err := proto.Unmarshal(paramsBuff, pbData); err != nil {
			return err
		}
	}

	buff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := proto.Unmarshal(buff, pbData); err != nil {
		return err
	}

	return nil
}
