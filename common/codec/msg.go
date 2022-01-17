/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : msg.go
*   coder: zemanzeng
*   date : 2022-01-16 11:43:35
*   desc :
*
================================================================*/

package codec

import "context"

type msg struct {
	logger  interface{}
	context context.Context
}

func newmsg() *msg {
	return &msg{}
}

func (msg *msg) Logger() interface{} {
	if msg == nil {
		return nil
	}
	return msg.logger
}

func (msg *msg) SetLogger(logger interface{}) {
	msg.logger = logger
}
