/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : service.go
*   coder: zemanzeng
*   date : 2022-02-13 10:58:07
*   desc : service body
*
================================================================*/

package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"

	pb "github.com/wenruo95/jerry/codec/proto"
	"google.golang.org/protobuf/proto"
)

type Body struct {
	Version     byte
	ServiceLen  uint32
	ServiceBody []byte
	DataLen     uint32
	Data        []byte
}

const (
	DefaultVersion = 0x1
)

func PackBody(serviceBody *pb.ServiceBody, dataMsg proto.Message) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := buf.WriteByte(DefaultVersion); err != nil {
		return nil, err
	}

	// service msg
	serviceBuff, err := proto.Marshal(serviceBody)
	if err != nil {
		return nil, err
	}

	binary.Write(buf, binary.BigEndian, len(serviceBuff))
	if _, err := buf.Write(serviceBuff); err != nil {
		return nil, err
	}

	// data msg
	dataBuff, err := proto.Marshal(dataMsg)
	if err != nil {
		return nil, err
	}

	binary.Write(buf, binary.BigEndian, len(dataBuff))
	if _, err := buf.Write(dataBuff); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnpackBody(buff []byte) (proto.Message, *pb.ServiceBody, error) {
	return UnpackBodyByReader(bytes.NewBuffer(buff))
}

func UnpackBodyByReader(reader io.Reader) (proto.Message, *pb.ServiceBody, error) {
	versionData := make([]byte, 1)
	if _, err := io.ReadFull(reader, versionData); err != nil {
		return nil, nil, err
	}

	switch versionData[0] {
	case 1:
		return UnpackV1(reader)
	default:
		return nil, nil, errors.New("not supported version:" + strconv.Itoa(int(versionData[0])))
	}
}

func UnpackV1(reader io.Reader) (proto.Message, *pb.ServiceBody, error) {
	var serviceLen uint32
	if err := binary.Read(reader, binary.BigEndian, &serviceLen); err != nil {
		return nil, nil, err
	}
	serviceBuff := make([]byte, serviceLen)
	if _, err := io.ReadFull(reader, serviceBuff); err != nil {
		return nil, nil, err
	}
	serviceBody := new(pb.ServiceBody)
	if err := proto.Unmarshal(serviceBuff, serviceBody); err != nil {
		return nil, nil, err
	}

	var dataLen uint32
	if err := binary.Read(reader, binary.BigEndian, &dataLen); err != nil {
		return nil, nil, err
	}
	dataBuff := make([]byte, dataLen)
	if _, err := io.ReadFull(reader, dataBuff); err != nil {
		return nil, nil, err
	}
	var dataBody proto.Message
	if err := proto.Unmarshal(dataBuff, dataBody); err != nil {
		return nil, nil, err
	}

	return dataBody, serviceBody, nil
}
