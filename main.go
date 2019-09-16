package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket/layers"
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
		fmt.Printf("Create client fail: %v \n", err)
	}

	defer client.Close()

	eventChan := make(chan *layers.SIP, 16)
	go client.Recv(eventChan)
	go client.ProcessPacket(eventChan)


	for {
		ch := time.After(5 * time.Second)
		select {
		case <-ch:
			fmt.Println("Try to register to srever...")
			if !client.Registered {
				client.Register()
			} else {
				fmt.Println("注册完成")
			}
		}
	}

}
