package main

import (
	"flag"
	"fmt"
	"iot-video-monitor/b"
	"iot-video-monitor/config"
	"time"
)

var (
	cfgFile = flag.String("config", `./examples/conf/sample.json`, "Configuration file")
)

func main() {
	fmt.Println("Start video monitor app...")

	cfg, err := config.NewConfig(*cfgFile)
	if err != nil {
		fmt.Printf("Parse config file fail: %v", err)
	}

	client, err := b.NewClient(cfg)
	if err != nil {
		fmt.Println("Create client fail: %v", err)
	}

	defer client.Close()

	go client.Recv()

	ch := time.After(5 * time.Second)
	for {
		select {
		case <-ch:
			fmt.Println("Try to register to srever...")
			client.Register()
			//default:
			//	fmt.Printf("Listening...\n")
			//read from UDPConn here
		}
	}

	//
	//
	//// write a message to server
	//packetData := map[string]string{"Call-ID": "306366781@172_16_254_66", "Contact": "<sip:bob@172.16.254.66:5060>"}
	//p := gopacket.NewPacket(packetData, LinkTypeEthernet, gopacket.Default)
	//
	//
	//
	//
	//
	//
	//n, err := conn.Write(message)
	//
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//// receive message from server
	//buffer := make([]byte, 1024)
	//n, addr, err := conn.ReadFromUDP(buffer)
	//
	//fmt.Println("UDP Server : ", addr)
	//fmt.Println("Received from UDP server : ", string(buffer[:n]))

}
