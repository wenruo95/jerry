/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : server.go
*   coder: zemanzeng
*   date : 2021-02-04 05:03:19
*   desc : server
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
	ServerStatusInit = iota
	ServerStatusServing
	ServerStatusListening
	ServerStatusStoped
)

var (
	ServerMinTimeoutMs     int64 = 100
	ServerDefaultTimeoutMs int64 = 30 * 1000
)

type ClientHandler interface {
	OnConnect(conn *ClientConn)
	OnMessage(conn *ClientConn, body []byte, messageFlag byte, txid uint32)
	OnDisconnect(conn *ClientConn, reason string)
}

type Server struct {
	addr     string
	status   int32
	handler  ClientHandler
	listener *net.TCPListener
	cancel   context.CancelFunc
}

func NewServer(addr string, handler ClientHandler) *Server {
	return &Server{
		addr:    addr,
		status:  ServerStatusInit,
		handler: handler,
	}
}

func (svr *Server) Addr() string {
	return svr.addr
}

func (svr *Server) Serve(timeout time.Duration) error {
	if svr == nil || len(svr.addr) == 0 || svr.handler == nil {
		return errors.New("invalid server args")
	}
	if !atomic.CompareAndSwapInt32(&svr.status, ServerStatusInit, ServerStatusServing) {
		return errors.New("invalid server status:" + strconv.FormatInt(int64(svr.status), 10))
	}

	addr, err := net.ResolveTCPAddr("tcp", svr.addr)
	if err != nil {
		return errors.New("resolve addr:" + svr.addr + " error:" + err.Error())
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return errors.New("listen addr:" + svr.addr + " error:" + err.Error())
	}
	svr.listener = listener

	// server has stoped
	if !atomic.CompareAndSwapInt32(&svr.status, ServerStatusServing, ServerStatusListening) {
		return errors.New("server has stoped. addr:" + svr.addr)
	}

	if timeout.Milliseconds() < ServerMinTimeoutMs {
		timeout = time.Duration(ServerDefaultTimeoutMs) * time.Millisecond
	}

	ctx, cancel := context.WithCancel(context.Background())
	svr.cancel = cancel
	for {
		select {
		case <-ctx.Done():
			return errors.New("context has canceled")

		default:
			conn, err := listener.AcceptTCP()
			if err != nil {
				return errors.New(svr.addr + " accept incomming connection error:" + err.Error())
			}
			go NewClientConn(conn, svr.handler).serve(ctx, timeout)
		}

	}
}

func (svr *Server) Close() error {
	if atomic.CompareAndSwapInt32(&svr.status, ServerStatusListening, ServerStatusStoped) {
		svr.cancel()
		return svr.listener.Close()
	}

	atomic.StoreInt32(&svr.status, ServerStatusStoped)
	return nil
}

type ClientConn struct {
	*net.TCPConn
	working int32
	closeCh chan struct{}
	handler ClientHandler
}

func NewClientConn(conn *net.TCPConn, handler ClientHandler) *ClientConn {
	return &ClientConn{
		TCPConn: conn,
		closeCh: make(chan struct{}, 2),
		working: 0,
		handler: handler,
	}
}

func (client *ClientConn) serve(ctx context.Context, timeout time.Duration) {
	client.handler.OnConnect(client)
	client.working = 1

	timer := time.NewTimer(timeout)
	for {
		select {
		case <-ctx.Done():
			client.stopAndCallDisconnect("context canceled")
			return

		case <-timer.C:
			client.stopAndCallDisconnect("conn exceed " + timeout.String())
			return

		case <-client.closeCh:
			close(client.closeCh)
			client.closeCh = nil
			return

		default:
			body, messageFlag, txid, err := Unpack(client.TCPConn)
			if err != nil {
				client.stopAndCallDisconnect("unpack " +
					client.RemoteAddr().String() + " msg error:" + err.Error())
				return
			}
			timer.Reset(timeout)

			client.handler.OnMessage(client, body, messageFlag, txid)
		}
	}

}

func (client *ClientConn) Send(body []byte, messageFlag byte, txid uint32) error {
	if atomic.LoadInt32(&client.working) == 0 {
		return errors.New("conn not working")
	}
	return PackWrite(client.TCPConn, messageFlag, txid, body)
}

func (client *ClientConn) stopAndCallDisconnect(reason string) {
	atomic.StoreInt32(&client.working, 0)

	if client.closeCh != nil {
		close(client.closeCh)
		client.closeCh = nil
	}

	if err := client.Close(); err != nil {
		s := fmt.Sprintf("conn(local:%s remote:%s) error:%s",
			client.LocalAddr().String(), client.RemoteAddr().String(), err.Error())
		reason = reason + " closed:" + s
	}
	client.handler.OnDisconnect(client, reason)
}

func (client *ClientConn) Stop() error {
	if !atomic.CompareAndSwapInt32(&client.working, 1, 0) {
		return errors.New("client has stoped")
	}

	client.closeCh <- struct{}{}
	return client.Close()
}
