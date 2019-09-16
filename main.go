package main

import (
	"flag"
	"fmt"
	"iot-video-monitor/b"
	"iot-video-monitor/config"
	"time"
	"github.com/marv2097/siprocket"
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

	eventChan := make(chan siprocket.SipMsg, 16)
	go client.Recv(eventChan)
	go client.ProcessPacket(eventChan)


	for {
		ch := time.After(5 * time.Second)
		select {
		case <-ch:
			if !client.Registered {
				client.Register()
			}
		}
	}

}
