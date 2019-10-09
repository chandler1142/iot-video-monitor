package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	dataBytes, err := ioutil.ReadFile("C:\\Users\\86188\\Downloads\\test.ts")

	packet := b.NewEmptyRTPPacket(33, 0, time.Now().Second())
	var i = 0
	count := 0
	for ; i < len(dataBytes); i += b.TS_PACKET_SIZE {
		data := dataBytes[i : i+b.TS_PACKET_SIZE]
		if data[0] != 0x47 {
			fmt.Println("Error packet...")
			continue
		}
		count += len(data)
		packet.SetData(data, false)

		if count+b.TS_PACKET_SIZE > b.MTU || i+b.TS_PACKET_SIZE >= len(dataBytes) {
			packet.SequenceIncrement()
			client.Conn.Write(packet.GetBytes())
			fmt.Printf("Send Bytes: %d \n", len(packet.GetBytes()))
			packet.SetData([]byte{}, true)
			count = 0
			time.Sleep(10 * time.Millisecond)
		}
	}


	////	[RTP-Header] Version: 2, Padding: 0, Extension: 0, CC: 0, Marker: 0, PayloadType: 26, SequenceNumber: 11, TimeStamp: 1100
	////801a000b0000044c00000539
	//packet := b.NewRTPPacket(96, 11, 1100, []byte{1, 2, 3}, 3)
	//fmt.Println(hex.EncodeToString(packet.Header))
	//
	//buf, err := ioutil.ReadFile("examples/video/11.mp4")
	//if err != nil {
	//	fmt.Println("error")
	//}
	//
	//totalLen := len(buf)
	//fmt.Println(totalLen)
	//
	//var times = 0
	//for i := 0; i < totalLen; {
	//	frameLen := 1024
	//
	//	fmt.Printf("frame length: %d \n", frameLen)
	//	data := buf[i : i+int(frameLen)]
	//	//rtpPacket := b.NewRTPPacket(26, times, times*100, data, int(frameLen))
	//	//packet := append(rtpPacket.Header, rtpPacket.Payload...)
	//
	//	for l := 0; l < len(data); l += 1024 {
	//		smallData := data[l : l+1024]
	//		rtpPacket := b.NewRTPPacket(26, times, times*100, smallData, len(smallData))
	//		packet := append(rtpPacket.Header, rtpPacket.Payload...)
	//		n, err := client.Conn.Write(packet)
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		fmt.Printf("Send : %d \n", n)
	//		time.Sleep(5 * time.Second)
	//	}
	//	times ++
	//	i = i + int(frameLen)
	//}

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
