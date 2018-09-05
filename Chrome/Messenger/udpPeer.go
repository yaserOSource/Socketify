package main

import (
	"net"
	"os"
)

var udpPeerSocket *net.UDPConn

func udpPeer(addrStr string) {
	address, error := net.ResolveUDPAddr("udp", addrStr)
	if error != nil {
		write(message{
			Event: "open",
			Error: "cannot resolve udp address",
			Debug: error.Error(),
		})
		os.Exit(1)
	}

	udpPeerSocket, error = net.ListenUDP("udp", address)
	if error != nil {
		write(message{
			Event: "open",
			Error: "cannot open udp socket",
			Debug: error.Error(),
		})
		os.Exit(1)
	}

	defer udpPeerSocket.Close()
	go udpPeerReceive()
	write(message{
		Event:   "open",
		Address: udpPeerSocket.LocalAddr().String(),
	})

	for {
		msg := read()
		switch msg.Event {
		case "error":
			write(message{
				Event: "close",
				Error: msg.Error,
				Debug: msg.Debug,
			})
			os.Exit(1)
			break
		case "send":
			udpPeerSend(msg)
			break
		case "close":
			write(message{
				Event: "close",
			})
			udpPeerSocket.Close()
			os.Exit(0)
			break
		}
	}
}

func udpPeerSend(msg message) {
	address, error := net.ResolveUDPAddr("udp", msg.Address)
	if error != nil {
		return
	}

	if address.IP == nil || address.Port == 0 {
		return
	}

	_, error = udpPeerSocket.WriteToUDP([]byte(msg.Payload), address)
	if error != nil {
		write(message{
			Event: "close",
			Error: "cannot write to udp socket",
			Debug: error.Error(),
		})
		os.Exit(1)
	}
}

func udpPeerReceive() {
	buffer := make([]byte, 1500)

	var length int
	var address *net.UDPAddr
	var error error

	for {
		length, address, error = udpPeerSocket.ReadFromUDP(buffer)
		if error != nil {
			write(message{
				Event: "close",
				Error: "cannot read from udp socket",
				Debug: error.Error(),
			})
			os.Exit(1)
		}

		write(message{
			Event:   "receive",
			Address: address.String(),
			Payload: string(buffer[:length]),
		})
	}
}