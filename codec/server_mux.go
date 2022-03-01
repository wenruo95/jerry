/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : server_mux.go
*   coder: zemanzeng
*   date : 2022-02-27 19:25:34
*   desc : server mux
*
================================================================*/

package codec

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"sync/atomic"

	pb "github.com/wenruo95/jerry/codec/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type ServiceHandler func(ctx context.Context, req proto.Message) (proto.Message, error)

var serverMux *ServerMux

type ServerMux struct {
	txid uint32
	muxs map[string]*MuxInfo
}

type MuxInfo struct {
	cmd     string
	msgName string
	handler ServiceHandler
}

func init() {
	serverMux = &ServerMux{
		txid: rand.Uint32(),
		muxs: make(map[string]*MuxInfo),
	}

}

func Register(cmd string, msgName string, handler ServiceHandler) {
	if _, exist := serverMux.muxs[cmd]; exist {
		panic("cmd:" + cmd + "has registered!!!")
	}

	serverMux.muxs[cmd] = &MuxInfo{cmd: cmd, msgName: msgName, handler: handler}
}

func GenTxid() uint32 {
	return atomic.AddUint32(&serverMux.txid, 1)
}

func NewReqTypeByCmd(cmd string) (proto.Message, error) {
	mux, exist := serverMux.muxs[cmd]
	if !exist {
		return nil, errors.New("cmd:" + cmd + " not exist")
	}

	fullName := protoreflect.FullName(mux.msgName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
	if err != nil {
		return nil, err
	}
	return msgType.New().Interface(), nil
}

func MuxHandler(ctx context.Context, w io.Writer, cmd string, buff []byte, metaData *pb.ServerMetaData) error {
	req, err := NewReqTypeByCmd(cmd)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(buff, req); err != nil {
		return err
	}

	rsp, err := serverMux.muxs[cmd].handler(ctx, req)
	if err != nil {
		return err
	}

	bodyBuff, err := proto.Marshal(rsp)
	if err != nil {
		return err
	}
	metaBuff, err := proto.Marshal(metaData)
	if err != nil {
		return err
	}

	txid := GenTxid()
	rspBuff, err := PackMsgV2(txid, bodyBuff, metaBuff)
	if err != nil {
		return err
	}

	if _, err := w.Write(rspBuff); err != nil {
		return err
	}
	return nil
}
