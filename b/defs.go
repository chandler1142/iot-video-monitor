package b

import (
	"net"
)

type Client struct {
	Conn                net.Conn
	CallId              string
	FromTag             string
	MessageTemplatePath string
	Registered          bool
	remoteAddr          *net.UDPAddr
	localAddr           *net.UDPAddr
	SysAddrCode         string
	UserAddrCode        string
}
