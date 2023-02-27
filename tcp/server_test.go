/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : server_test.go
*   coder: zemanzeng
*   date : 2021-02-04 05:03:56
*   desc : tcp测试用例
*
================================================================*/

package tcp

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type TestTCPServerHandler struct {
}

func (handler *TestTCPServerHandler) OnConnect(conn *ClientConn) {
	log.Printf("[I] server: new client connection. addr:%v", conn.RemoteAddr().String())
}

func (handler *TestTCPServerHandler) OnMessage(conn *ClientConn, body []byte, messageFlag byte, txid uint32) {
	log.Printf("[I] server: recieve message. len:%v flag:%v txid:%v body:%s", len(body), messageFlag, txid, string(body))
	conn.Send([]byte("server_hello_world_"+strconv.FormatUint(uint64(txid), 10)), messageFlag, txid)
}

func (handler *TestTCPServerHandler) OnDisconnect(conn *ClientConn, reason string) {
	log.Printf("[I] server: client disconnected. addr:%v reason:%v", conn.RemoteAddr().String(), reason)
}

func Test_ServerTimeout(t *testing.T) {
	svr := NewServer("127.0.0.1:8000", &TestTCPServerHandler{})
	go svr.Serve(10 * time.Second)
	time.Sleep(time.Second)
	svr.Close()
}
