package b

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"iot-video-monitor/config"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type Client struct {
	conn                net.Conn
	CallId              string
	FromTag             string
	MessageTemplatePath string
	registered          bool
	remoteAddr          *net.UDPAddr
	localAddr           *net.UDPAddr
	cfg                 *config.Config
}

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
		registered:          false,
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

//REGISTER sip:22.46.93.183 SIP/2.0
//From: <sip:100010000003010002@22.46.93.183>;tag=447226015
//To: <sip:100010000003010002@22.46.93.183>
//Call-ID: 1702005310
//Via: SIP/2.0/UDP 22.46.93.196:5060;rport;branch=z9hG4bK916029210
//CSeq: 38 REGISTER
//Contact: <sip:100010000003010002@22.46.93.196:5060>
//Max-Forwards: 70
//User-Agent: IP CAMERA
//Expires: 100
//Content-Length: 0
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
func (client *Client) Recv() {
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
		fmt.Printf("\n===== Receive message: =====\n%v \n", string(buf[:n]))
	}
}

func GetRandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
