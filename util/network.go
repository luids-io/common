// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package util

import (
	"fmt"
	"net"
	"strings"
)

// ParseListenURI returns proto, addrs and error if URI is not valid
func ParseListenURI(s string) (proto string, addr string, err error) {
	err = nil
	if strings.HasPrefix(s, "unix://") {
		proto = "unix"
		addr = s[7:]
	} else if strings.HasPrefix(s, "tcp://") {
		proto = "tcp"
		addr = s[6:]
	} else {
		err = fmt.Errorf("invalid prefix in %v", s)
	}
	return
}

// Listener returns a listener socket from an uri
func Listener(uri string) (net.Listener, error) {
	proto, addr, err := ParseListenURI(uri)
	if err != nil {
		return nil, fmt.Errorf("cannot parse address '%v': %v", uri, err)
	}
	lis, err := net.Listen(proto, addr)
	if err != nil {
		return nil, fmt.Errorf("cannot listen socket '%v': %v", uri, err)
	}
	return lis, nil
}
