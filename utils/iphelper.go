/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : iphelper.go
*   coder: zemanzeng
*   date : 2021-09-27 18:28:24
*   desc : ip helper
*
================================================================*/

package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func InetNtoa(ip uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xff)
	bytes[1] = byte((ip >> 8) & 0xff)
	bytes[2] = byte((ip >> 16) & 0xff)
	bytes[3] = byte((ip >> 24) & 0xff)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func IPtoI(ip string) uint32 {
	ips := net.ParseIP(ip)
	if ips == nil {
		return 0
	}

	if len(ips) == 16 {
		return binary.BigEndian.Uint32(ips[12:16])
	}
	return binary.BigEndian.Uint32(ips)
}

func ItoIP(ip uint32) string {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, ip); err != nil {
		return ""
	}

	b := buf.Bytes()
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func InetAton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0]) // nolint
	b1, _ := strconv.Atoi(bits[1]) // nolint
	b2, _ := strconv.Atoi(bits[2]) // nolint
	b3, _ := strconv.Atoi(bits[3]) // nolint

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

func GetInterIP(name string) (*string, error) {
	interfaces, e := net.Interfaces()
	if e != nil {
		return nil, e
	}
	for _, inter := range interfaces {
		if inter.Name == name {
			addrs, e := inter.Addrs()
			if e != nil {
				return nil, e
			}
			if 0 == len(addrs) {
				return nil, errors.New("empty interface")
			}
			valid := regexp.MustCompile("[0-9.]+")
			fileds := valid.FindAllString(addrs[0].String(), -1)
			return &fileds[0], nil
		}
	}
	return nil, errors.New("invalid interface")
}

func GetLocalIp() string {
	for i := 0; i < 10; i++ {
		strInf := fmt.Sprintf("eth%d", i)
		ip, err := GetInterIP(strInf)
		if err == nil && *ip != "" {
			return *ip
		}
	}

	return "127.0.0.1"
}

var localIp = GetLocalIp()

func GetFastLocalIp() string {
	return localIp
}

func ParseAddr(addr string) (ip string, port uint16, err error) {
	addrs := strings.Split(addr, ":")
	if len(addrs) != 2 {
		return "", 0, errors.New("不合法的地址格式")
	}

	if ips := net.ParseIP(addrs[0]); ips == nil {
		return "", 0, errors.New("表示不合法的ip地址格式")
	}

	portTmp, err := strconv.Atoi(addrs[1])
	if err != nil || portTmp <= 0 || portTmp >= 65536 {
		return "", 0, errors.New("port not valid")
	}
	return addrs[0], uint16(portTmp), nil
}
