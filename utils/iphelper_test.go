/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : iphelper_test.go
*   coder: zemanzeng
*   date : 2021-09-27 18:30:06
*   desc :
*
================================================================*/

package utils

import (
	"fmt"
	"net"
	"testing"
)

func TestGetInterIP(t *testing.T) {
	ip_str := string("eth0")
	ip, err := GetInterIP(ip_str)
	if err != nil {
		t.Error("GetInterIP failed")
	}
	fmt.Println(ip)
}

func TestInetNtoa(t *testing.T) {
	ip := net.IPv4(127, 0, 0, 1)
	ip_num := InetAton(ip)
	if ip_num != 0 {
		t.Error("InetNtoa failed")
	}
}

func TestInetAton(t *testing.T) {
	net.IPv4(0, 0, 0, 0)
	InetNtoa(0)
}

func TestGetLocalIp(t *testing.T) {
	t.Logf("ip:%s", GetLocalIp())
}
