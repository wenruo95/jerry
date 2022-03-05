/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : msg.go
*   coder: zemanzeng
*   date : 2022-01-16 11:43:35
*   desc : Msg实现
*
================================================================*/

package codec

import "context"

type msg struct {
	context context.Context
	logger  interface{}

	envid   uint32 // 0-正式环境 其余均为测试环境
	uin     uint64 // 用户uin
	appid   uint32 // 服务appid
	event   string // 服务的api
	traceid string // 请求的traceid
}

func newmsg() *msg {
	return &msg{}
}

func (msg *msg) Context() context.Context {
	return msg.context
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

func (msg *msg) EnvId() uint32 {
	return msg.envid
}

func (msg *msg) SetEnvId(envid uint32) {
	msg.envid = envid
}

func (msg *msg) Uin() uint64 {
	return msg.uin
}

func (msg *msg) SetUin(uin uint64) {
	msg.uin = uin
}

func (msg *msg) AppId() uint32 {
	return msg.appid
}

func (msg *msg) SetAppId(appid uint32) {
	msg.appid = appid
}

func (msg *msg) Event() string {
	return msg.event
}

func (msg *msg) SetEvent(event string) {
	msg.event = event
}

func (msg *msg) TraceId() string {
	return msg.traceid
}

func (msg *msg) SetTraceId(traceId string) {
	msg.traceid = traceId
}
