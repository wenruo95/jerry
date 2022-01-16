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

type CtxMsg interface {
	Logger() Logger
	SetLogger(logger Logger)
}

func MustCtxMsg(ctx context.Context) (context.Context, CtxMsg) {
	ctxMsgI := ctx.Value(ctxKey)
	if ctxMsg, ok := ctxMsgI.(CtxMsg); ok {
		return ctx, ctxMsg
	}

	ctxMsg := NewCtxMeta()
	newCtx := context.WithValue(ctx, ctxKey, ctxMsg)
	return newCtx, ctxMsg
}
