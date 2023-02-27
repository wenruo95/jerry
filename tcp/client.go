/*================================================================

*   Copyright (C) 2021. All rights reserved.
*
*   file : client.go
*   coder: zemanzeng
*   date : 2021-02-04 05:04:45
*   desc : tcp client
*
================================================================*/

package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	ClientStatusInit = iota
	ClientStatusServing
	ClientStatusWorking
	ClientStatusStoped
)

var (
	ClientMinTimeoutMs     int64 = 100
	ClientDefaultTimeoutMs int64 = 30 * 1000
)

type ServerHandler interface {
	OnConnect(cli *Client)
	OnMessage(cli *Client, body []byte, messageFlag byte, txid uint32)
	OnDisconnect(cli *Client, reason string)
}

type Client struct {
	*net.TCPConn
	addr    string
	status  int32
	handler ServerHandler
	closeCh chan struct{}
}

func NewClient(addr string, handler ServerHandler) *Client {
	return &Client{
		addr:    addr,
		status:  ClientStatusInit,
		handler: handler,
		closeCh: make(chan struct{}),
	}
}

func (cli *Client) Serve(ctx context.Context, timeout time.Duration) error {
	if cli == nil || len(cli.addr) == 0 || cli.handler == nil {
		return errors.New("invalid client args")
	}
	if !atomic.CompareAndSwapInt32(&cli.status, ClientStatusInit, ClientStatusServing) {
		return errors.New("invalid status:" + strconv.FormatInt(int64(cli.status), 10))
	}

	addr, err := net.ResolveTCPAddr("tcp", cli.addr)
	if err != nil {
		return errors.New("resolve addr:" + cli.addr + " error:" + err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return errors.New("dial server:" + cli.addr + " error:" + err.Error())
	}
	cli.TCPConn = conn
	cli.handler.OnConnect(cli)

	// stoped
	if !atomic.CompareAndSwapInt32(&cli.status, ClientStatusServing, ClientStatusWorking) {
		return errors.New("invalid status:" + strconv.FormatInt(int64(cli.status), 10))
	}

	if timeout.Milliseconds() < ClientMinTimeoutMs {
		timeout = time.Duration(ClientDefaultTimeoutMs) * time.Millisecond
	}

	timer := time.NewTimer(timeout)
	for {
		select {
		case <-ctx.Done():
			cli.stopAndCallDisconnect("context canceled")
			return nil

		case <-timer.C:
			cli.stopAndCallDisconnect("conn exceed " + timeout.String())
			return nil

		case <-cli.closeCh:
			close(cli.closeCh)
			cli.closeCh = nil
			return nil

		default:
			body, messageFlag, txid, err := Unpack(conn)
			if err != nil {
				return errors.New(conn.LocalAddr().String() + " unpack " +
					conn.RemoteAddr().String() + " msg error:" + err.Error())
			}
			timer.Reset(timeout)

			cli.handler.OnMessage(cli, body, messageFlag, txid)
		}
	}

	return nil
}

func (cli *Client) Send(body []byte, messageFlag byte, txid uint32) error {
	if atomic.LoadInt32(&cli.status) != ClientStatusWorking {
		return errors.New("conn not working")
	}

	return PackWrite(cli.TCPConn, messageFlag, txid, body)
}

func (cli *Client) stopAndCallDisconnect(reason string) {
	atomic.StoreInt32(&cli.status, ClientStatusStoped)

	close(cli.closeCh)
	cli.closeCh = nil

	if err := cli.Close(); err != nil {
		s := fmt.Sprintf("conn(local:%s remote:%s) error:%s",
			cli.LocalAddr().String(), cli.RemoteAddr().String(), err.Error())
		reason = reason + " closed:" + s
	}
	cli.handler.OnDisconnect(cli, reason)
}

func (cli *Client) Stop() error {
	if atomic.CompareAndSwapInt32(&cli.status, ClientStatusWorking, ClientStatusStoped) {
		cli.closeCh <- struct{}{}
		return cli.Close()
	}

	atomic.StoreInt32(&cli.status, ClientStatusStoped)
	return nil
}
