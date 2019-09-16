package b

import (
	"encoding/hex"
	"fmt"
	"github.com/google/gopacket/layers"
	"io/ioutil"
	"iot-video-monitor/config"
	"log"
	"net"
	"strconv"
)

func NewClient(cfg *config.Config) (*Client, error) {

	remoteAddr := cfg.ServerIp + ":" + strconv.Itoa(cfg.ServerPort)
	RemoteAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		fmt.Printf("Remote Server address error: %v", err)
		return nil, err
	}

	localAddr := cfg.LocalIp + ":" + strconv.Itoa(cfg.LocalPort)
	LocalAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		fmt.Printf("Local Server address error: %v", err)
		return nil, err
	}

	conn, err := net.DialUDP("udp", LocalAddr, RemoteAddr)
	if err != nil {
		fmt.Printf("Create UDP connection error: %v", err)
		return nil, err
	}

	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	return &Client{
		conn:                conn,
		MessageTemplatePath: cfg.TemplatePath,
		Registered:          false,
		remoteAddr:          RemoteAddr,
		localAddr:           LocalAddr,
		cfg:                 cfg,
		CallId:              GetRandomString(8),
		FromTag:             GetRandomString(8),
	}, nil

}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Register() bool {
	buf, err := ioutil.ReadFile(client.MessageTemplatePath + "register")
	if err != nil {
		fmt.Printf("open file error: %v \n", err)
		return false
	}

	template := string(buf)
	message := fmt.Sprintf(template,
		client.remoteAddr.IP,
		client.cfg.SysAddrCode,
		client.remoteAddr.IP,
		GetRandomString(8),
		client.cfg.SysAddrCode,
		client.remoteAddr.IP,
		client.CallId,
		client.localAddr.IP,
		client.localAddr.Port,
		client.cfg.SysAddrCode,
		client.localAddr.IP,
		client.localAddr.Port,
		3600,
	)
	client.conn.Write([]byte(message))
	fmt.Printf("\n ===== Send message： =====\n%s\n", message)
	return true
}

//启动监听终端接收指令的goroutine
func (client *Client) Recv(packetChan chan *layers.SIP) {
	if client.conn == nil {
		fmt.Println("robot connection has not initialized...")
		return
	}
	buf := make([]byte, 1024)

	//接收数据包，并解析
	for {
		n, err := client.conn.Read(buf[:])

		if n == 0 {
			continue
		}
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			fmt.Printf("received error packet is: %s \n", hex.EncodeToString(buf[:]))
			break
		}
		sipPacket := layers.NewSIP()
		sipPacket.DecodeFromBytes(buf[:n], nil)
		fmt.Printf("\n===== Receive message: =====\n%v \n", string(buf[:n]))
		packetChan <- sipPacket
	}
}

func (client *Client) ProcessPacket(packetChan chan *layers.SIP) {
	for {
		select {
		case packet := <-packetChan:
			if packet.GetFirstHeader("cseq") == "REGISTER" && packet.ResponseCode == 200 {
				client.Registered = true
			}


		}
	}
}
