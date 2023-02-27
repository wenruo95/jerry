/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : client_test.go
*   coder: zemanzeng
*   date : 2021-02-04 05:18:31
*   desc :
*
================================================================*/

package tcp

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"
)

type TestTCPClientHandler struct {
}

func (handler *TestTCPClientHandler) OnConnect(cli *Client) {
	log.Printf("[I] client: connect to server.")
}

func (handler *TestTCPClientHandler) OnMessage(cli *Client, body []byte, messageFlag byte, txid uint32) {
	log.Printf("[I] client: recieve message. len:%v flag:%v txid:%v body:%s", len(body), messageFlag, txid, string(body))
}

func (handler *TestTCPClientHandler) OnDisconnect(cli *Client, reason string) {
	log.Printf("[I] client: disconnect to server. reason:%v", reason)
}

func Test_Client(t *testing.T) {
	svr := NewServer("127.0.0.1:8000", &TestTCPServerHandler{})
	defer svr.Close()
	go func() {
		if err := svr.Serve(3 * time.Second); err != nil {
			t.Errorf("start server error:%v", err)
		}
	}()
	time.Sleep(time.Second)

	cli := NewClient("127.0.0.1:8000", &TestTCPClientHandler{})
	defer cli.Stop()

	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 50; i++ {
			txid := uint32(i + 100)
			cli.Send([]byte("client_hello_world_"+strconv.FormatUint(uint64(txid), 10)), 1, uint32(txid))
			time.Sleep(time.Second)
		}
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
		defer cancel()

		if err := cli.Serve(ctx, time.Duration(3)*time.Second); err != nil {
			t.Errorf("client serve error:%v", err)
		}
	}()

	time.Sleep(30 * time.Second)

}
