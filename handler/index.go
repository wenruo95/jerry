/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : index.go
*   coder: zemanzeng
*   date : 2022-01-16 16:30:03
*   desc :
*
================================================================*/

package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type DefaultHandlerReq struct {
}

type DefaultHandlerResp struct {
	Code int
	Msg  string
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	rsp := new(DefaultHandlerResp)
	defer func() {
		rspBody, err := json.Marshal(w)
		if err != nil {
			w.Write([]byte("internal error"))
			log.Printf("[ERROR] unmarshal rsp body error:" + err.Error())
			return
		}
		if _, err := w.Write(rspBody); err != nil {
			log.Printf("[ERROR] write rsp body error:" + err.Error())
			return
		}

		return
	}()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rsp.Code, rsp.Msg = 100, "read http body error"
		log.Printf("[ERROR] read body error:" + err.Error())
		return
	}

	req := new(DefaultHandlerReq)
	if err := json.Unmarshal(reqBody, req); err != nil {
		rsp.Code, rsp.Msg = 101, "unmarshal http body error"
		log.Printf("[ERROR] unmarshal body:%s error:%v", string(reqBody), err)
		return
	}

	// do some logic

}
