package b

import (
	"encoding/hex"
	"fmt"
	"github.com/marv2097/siprocket"
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
		CallId:              GetRandomString(8),
		FromTag:             GetRandomString(8),
		SysAddrCode:         cfg.SysAddrCode,
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
		client.SysAddrCode,
		client.remoteAddr.IP,
		client.FromTag,
		client.SysAddrCode,
		client.remoteAddr.IP,
		client.CallId,
		client.localAddr.IP,
		client.localAddr.Port,
		client.SysAddrCode,
		client.localAddr.IP,
		client.localAddr.Port,
		3600,
	)
	client.conn.Write([]byte(message))
	fmt.Printf("\n ===== Send message： =====\n%s\n", message)
	return true
}

//SIP/2.0 100 Trying
//Via: SIP/2.0/UDP 22.46.93.183:5060;rport=5060;branch=z9hG4bK492ad
//From: <sip:a978bab1a7c84c34bc3fd51d87c5c5da@22.46.93.183>;tag=40416f28
//To: <sip:100010000003010002@22.46.93.183>
//Call-ID: 2264a412a386@22.46.93.183
//CSeq: 1 INVITE
//User-Agent: IP CAMERA
//Content-Length: 0
func (client *Client) Trying() {
	buf, err := ioutil.ReadFile(client.MessageTemplatePath + "trying")
	if err != nil {
		fmt.Printf("open file error: %v \n", err)
		return
	}

	template := string(buf)
	message := fmt.Sprintf(template,
		client.remoteAddr.IP,
		client.remoteAddr.Port,
		client.remoteAddr.Port,
		client.UserAddrCode,
		client.remoteAddr.IP,
		client.FromTag,
		client.SysAddrCode,
		client.remoteAddr.IP,
		client.CallId,
	)
	client.conn.Write([]byte(message))
	fmt.Printf("\n ===== Send message： =====\n%s\n", message)
}

//启动监听终端接收指令的goroutine
func (client *Client) Recv(packetChan chan siprocket.SipMsg) {
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

		sipMsg := siprocket.Parse(buf)
		fmt.Printf("\n===== Receive message: =====\n%v \n", string(buf[:n]))
		packetChan <- sipMsg
	}
}

func (client *Client) ProcessPacket(packetChan chan siprocket.SipMsg) {
	for {
		select {
		case packet := <-packetChan:
			//注册成功的消息
			if string(packet.Cseq.Method) == "REGISTER" && string(packet.Req.StatusCode) == "200" {
				fmt.Println("注册完成")
				client.Registered = true
			}
			//调阅视频的消息
			if string(packet.Cseq.Method) == "INVITE" {
				fmt.Println("收到调阅视频请求，发送Trying")
				client.UserAddrCode = string(packet.From.User)
				client.CallId = string(packet.CallId.Value)
				client.Trying()
			}
		}
	}
}
