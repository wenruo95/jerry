/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : server_test.go
*   coder: zemanzeng
*   date : 2022-02-27 19:57:52
*   desc :
*
================================================================*/

package codec

import (
	"context"
	"testing"

	pb "github.com/wenruo95/jerry/codec/proto"
	"google.golang.org/protobuf/proto"
)

func TestServerCodec(t *testing.T) {

	fn := func(ctx context.Context, req proto.Message) (proto.Message, error) {
		return nil, nil
	}

	fullName := (&pb.ServerMetaData{}).ProtoReflect().Descriptor().FullName()
	Register("hello_world", string(fullName), fn)

	msg, err := NewReqTypeByCmd("hello_world")
	if err != nil {
		t.Errorf("new_req_type_by_cmd error:" + err.Error())
	}
	if _, ok := msg.(*pb.ServerMetaData); !ok {
		t.Errorf("msg interface trans error. msg:%T expect:%v", msg, fullName)
	}
	t.Logf("msg interface. msg:%T expect:%v", msg, fullName)

}
