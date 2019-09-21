package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"iot-video-monitor/b"
	"iot-video-monitor/config"
	"strconv"
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

	//	[RTP-Header] Version: 2, Padding: 0, Extension: 0, CC: 0, Marker: 0, PayloadType: 26, SequenceNumber: 11, TimeStamp: 1100
	//801a000b0000044c00000539
	packet := b.NewRTPPacket(26, 11, 1100, []byte{1, 2, 3}, 3)
	fmt.Println(hex.EncodeToString(packet.Header))

	buf, err := ioutil.ReadFile("movie.Mjpeg")
	if err != nil {
		fmt.Println("error")
	}

	fmt.Println(strconv.ParseInt(string(buf[:5]), 10, 64))

	fmt.Println(len(buf[5:]))

	totalLen := len(buf)
	fmt.Println(totalLen)

	var times = 0
	for i := 0; i < totalLen; {
		frameLen, err := strconv.ParseInt(string(buf[i:i+5]), 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("frame length: %d \n", frameLen)
		data := buf[i+5 : i+5+int(frameLen)]
		//rtpPacket := b.NewRTPPacket(26, times, times*100, data, int(frameLen))
		//packet := append(rtpPacket.Header, rtpPacket.Payload...)

		for l := 0; l < len(data); l += 1024 {
			smallData :=  data[l : l+1024]
			rtpPacket := b.NewRTPPacket(26, times, times*100, smallData, len(smallData))
			packet := append(rtpPacket.Header, rtpPacket.Payload...)
			n, err := client.Conn.Write(packet)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Send : %d \n", n)
			time.Sleep(1 * time.Second)
		}
		times ++
		i = i + int(frameLen) + 5
	}

	//
	//eventChan := make(chan siprocket.SipMsg, 16)
	//go client.Recv(eventChan)
	//go client.ProcessPacket(eventChan)
	//
	//
	//for {
	//	ch := time.After(5 * time.Second)
	//	select {
	//	case <-ch:
	//		if !client.Registered {
	//			client.Register()
	//		}
	//	}
	//}

}
