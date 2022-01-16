/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : ctxmeta.go
*   coder: zemanzeng
*   date : 2022-01-16 11:43:35
*   desc :
*
================================================================*/

package codec

type CtxMeta struct {
	logger Logger
}

func NewCtxMeta() *CtxMeta {
	return &CtxMeta{}
}

func (meta *CtxMeta) Logger() Logger {
	return meta.logger
}

func (meta *CtxMeta) SetLogger(logger Logger) {
	meta.logger = logger
}
