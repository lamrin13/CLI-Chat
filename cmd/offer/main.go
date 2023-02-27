package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	b64 "encoding/base64"

	"github.com/pion/webrtc/v2"
)

func main() {

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)

	if err != nil {
		fmt.Println("error creating new connection ", err)
		os.Exit(1)
	}

	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		fmt.Println(pcs.String())
	})

	// Creating offer
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		fmt.Println("error creating new offer ", err)
		os.Exit(1)
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		fmt.Println("error while setting local description ", err)
		os.Exit(1)
	}

	fmt.Println(b64.StdEncoding.EncodeToString([]byte(offer.SDP)))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter sdp: ")
	sdp, _ := reader.ReadString('\n')
	decodedSDP, _ := b64.StdEncoding.DecodeString(sdp)

	answer := webrtc.SessionDescription{
		Type: 3,
		SDP:  string(decodedSDP),
	}
	err = pc.SetRemoteDescription(answer)
	if err != nil {
		fmt.Println("error while setting remote description ", err)
		os.Exit(1)
	}

	data, err := pc.CreateDataChannel("Nirmal", nil)
	if err != nil {
		fmt.Println("error while creating data channel ", err)
		os.Exit(1)
	}
	data.OnOpen(func() {
		err = data.SendText("Hello from offer!")
		if err != nil {
			fmt.Println("error sending message from offer ", err)
			os.Exit(1)
		}
	})

	data.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from %s: %s", data.Label(), string(msg.Data))
	})

	message := "a"
	pc.OnDataChannel(func(data *webrtc.DataChannel) {
		data.OnOpen(func() {
			fmt.Printf("New data channel %s\n", data.Label())
		})

		data.OnMessage(func(msg webrtc.DataChannelMessage) {
			if message == "" {
				fmt.Printf("\r")
			}
			fmt.Printf("Message from %s: %s", data.Label(), string(msg.Data))
			fmt.Print("\nMessage to send: ")
		})
	})
	for {
		fmt.Print("Message to send: ")
		message, _ = reader.ReadString('\n')
		data.SendText(strings.TrimSuffix(message, "\n"))
		message = ""
	}
}
