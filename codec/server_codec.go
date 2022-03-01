/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : server_codec.go
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
	"fmt"
	"io"
)

const (
	DefaultHead byte = 0x33
	DefaultTail byte = 0x35

	MsgVersionV1 byte = 0x01
	MsgVersionV2 byte = 0x02
)

// head(1) | version(1) | txid(4) | flag(4) | bodyLen(4) | bodyData(m) | tail(1)
func PackMsgV1(txid uint32, flag uint32, bodyData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := buf.WriteByte(DefaultHead); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(MsgVersionV1); err != nil {
		return nil, err
	}
	binary.Write(buf, binary.BigEndian, txid)
	binary.Write(buf, binary.BigEndian, flag)
	binary.Write(buf, binary.BigEndian, len(bodyData))
	if _, err := buf.Write(bodyData); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func UnpackMsgV1(reader io.Reader) (txid uint32, flag uint32, bodyData []byte, err error) {
	pre := make([]byte, 2)
	if _, err = io.ReadFull(reader, pre); err != nil {
		return
	}
	if head := pre[0]; head != DefaultHead {
		err = fmt.Errorf("head not matched %v!=%v", head, DefaultHead)
		return
	}
	if version := pre[1]; version != MsgVersionV1 {
		err = fmt.Errorf("version not matched %v!=%v", version, MsgVersionV1)
		return
	}

	if err = binary.Read(reader, binary.BigEndian, &txid); err != nil {
		return
	}
	if err = binary.Read(reader, binary.BigEndian, &flag); err != nil {
		return
	}

	var bodyLen uint32
	if err = binary.Read(reader, binary.BigEndian, &bodyLen); err != nil {
		return
	}
	if bodyLen == 0 {
		err = errors.New("body length is zero")
		return
	}

	leftBuff := make([]byte, bodyLen+1)
	if _, err = io.ReadFull(reader, bodyData); err != nil {
		return
	}
	if tail := leftBuff[len(leftBuff)-1]; tail != DefaultTail {
		err = fmt.Errorf("tail not matched %v!=%v", tail, DefaultTail)
		return
	}

	bodyData = leftBuff[:len(leftBuff)-1]
	return
}

// head(1) | version(1) | txid(4) | bodyLen(4) | metaLen(4) | bodyData(m) | metaData(n) | tail(1)
func PackMsgV2(txid uint32, bodyData []byte, metaData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := buf.WriteByte(DefaultHead); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(MsgVersionV1); err != nil {
		return nil, err
	}
	binary.Write(buf, binary.BigEndian, txid)
	binary.Write(buf, binary.BigEndian, len(bodyData))
	binary.Write(buf, binary.BigEndian, len(metaData))
	if _, err := buf.Write(bodyData); err != nil {
		return nil, err
	}
	if _, err := buf.Write(metaData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnpackMsgV2(reader io.Reader) (txid uint32, bodyData []byte, metaData []byte, err error) {
	pre := make([]byte, 2)
	if _, err = io.ReadFull(reader, pre); err != nil {
		return
	}
	if head := pre[0]; head != DefaultHead {
		err = fmt.Errorf("head not matched %v!=%v", head, DefaultHead)
		return
	}
	if version := pre[1]; version != MsgVersionV1 {
		err = fmt.Errorf("version not matched %v!=%v", version, MsgVersionV1)
		return
	}

	if err = binary.Read(reader, binary.BigEndian, &txid); err != nil {
		return
	}

	var bodyLen uint32
	if err = binary.Read(reader, binary.BigEndian, &bodyLen); err != nil {
		return
	}
	if bodyLen == 0 {
		err = errors.New("body length is zero")
		return
	}

	var metaLen uint32
	if err = binary.Read(reader, binary.BigEndian, &metaLen); err != nil {
		return
	}
	if metaLen == 0 {
		err = errors.New("meta length is zero")
		return
	}

	leftBuff := make([]byte, bodyLen+metaLen+1)
	if _, err = io.ReadFull(reader, bodyData); err != nil {
		return
	}
	if tail := leftBuff[len(leftBuff)-1]; tail != DefaultTail {
		err = fmt.Errorf("tail not matched %v!=%v", tail, DefaultTail)
		return
	}

	bodyData = leftBuff[:bodyLen]
	metaData = leftBuff[bodyLen : len(leftBuff)-1]

	return
}
