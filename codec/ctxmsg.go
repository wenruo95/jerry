/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : ctxmsg.go
*   coder: zemanzeng
*   date : 2022-01-16 11:22:28
*   desc : ctx msg
*
================================================================*/

package codec

import (
	"context"
)

type ContextMsgKey string

var ctxKey ContextMsgKey = ContextMsgKey("ctx_msg_key")

type Msg interface {
	Context() context.Context

	Logger() interface{}
	SetLogger(logger interface{})

	EnvId() uint32
	SetEnvId(envid uint32)

	Uin() uint64
	SetUin(uin uint64)

	AppId() uint32
	SetAppId(appid uint32)

	TraceId() string
	SetTraceId(traceId string)
}

func NewMsg(ctx context.Context) (context.Context, Msg) {
	msg := newmsg()
	newCtx := context.WithValue(ctx, ctxKey, msg)
	msg.context = newCtx
	return newCtx, msg
}

func EnsureMsg(ctx context.Context) (context.Context, Msg) {
	msgI := ctx.Value(ctxKey)
	if msg, ok := msgI.(*msg); ok {
		return ctx, msg
	}

	return NewMsg(ctx)
}

func Message(ctx context.Context) Msg {
	if msg, ok := ctx.Value(ctxKey).(*msg); ok {
		return msg
	}
	return &msg{context: ctx}
}
