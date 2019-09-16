package b

import (
	"iot-video-monitor/config"
	"net"
)

type Client struct {
	conn                net.Conn
	CallId              string
	FromTag             string
	MessageTemplatePath string
	Registered          bool
	remoteAddr          *net.UDPAddr
	localAddr           *net.UDPAddr
	cfg                 *config.Config
}
