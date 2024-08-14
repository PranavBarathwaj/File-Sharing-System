package udp

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/File-Sharing-System/server"
)

var (
	ClientConn        *net.UDPConn
	ServerConn        *net.UDPConn
	err               error
	ClientPort        int    = 9087
	ServerPort        int    = 9088
	BroadcastAddr     string = "255.255.255.255"
	createdClientConn bool   = false
	WaitingTime       int    = 1
)

type ReceivedMessage struct {
	Msg string
	Ip  string
}

func init() {

}

func CreateClientConn() {
	if createdClientConn {
		return
	}
	ClientConn, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: ClientPort,
		Zone: "",
	})
	if err != nil {
		panic("error in creating client connection " + err.Error())
	}
	log.Println("created Client connection of port", ClientPort)
	createdClientConn = true
}

func CloseClientConn() {
	if !createdClientConn {
		return
	}
	if err != ClientConn.Close() {
		log.Println("unable to close client connection")
	}
	createdClientConn = false
	log.Println("closed client conn")
}

func CreateServerConn() {
	ServerConn, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: ServerPort,
		Zone: "",
	})
	if err != nil {
		panic("error in creating server connection " + err.Error())
	}
	log.Println("created Server connection of port", ServerPort)
}

func CloseSeverConn() {
	if err != ServerConn.Close() {
		log.Println("unable to close server connection")
	}
	log.Println("closed server connction")
}

// sends message via any port to given addr and given port
func SendMsg(msg string, ip string, port int) {
	log.Println("trying to send", msg, "ip", ip, "port", port)
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		log.Println("parsing ip returned <nil> ")
		return
	}
	if msg == "" {
		log.Println("msg is empty")
		return
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ipAddr,
		Port: port,
		Zone: "",
	})
	if err != nil {
		log.Println("Unable To send Msg ", msg, "via", ipAddr.String(), "with port", port, err.Error())
		conn.Close()
		if err != nil {
			log.Println("error in closing the write connection " + err.Error())
		}
		return
	}
	log.Println("created a connection for sending message to ip:", ipAddr.String(), "via port", port)

	conn.Write([]byte(msg))
	log.Println("sent", msg, "from", conn.LocalAddr().String(), "to", conn.RemoteAddr().String())
	err = conn.Close()
	if err != nil {
		log.Println("error in closing the write connection " + err.Error())
	}
	log.Println("closed a created connection for sending msg")
}

// writes the message into channel
func ReadMsg(conn *net.UDPConn, channel chan ReceivedMessage) {
	log.Println("trying to read msg from", conn.LocalAddr().String())
	by := make([]byte, 14)
	n, addr, err := conn.ReadFrom(by)
	if err != nil {
		log.Println("error in reading the message " + err.Error())
		return
	}
	log.Println("read", n, "bytes", " from", addr.String(), "of value", string(by))
	s := ReceivedMessage{
		Msg: string(bytes.Trim(by, "\x00")),
		Ip:  strings.Split(addr.String(), ":")[0],
	}
	channel <- s
	log.Println("Read []byte form msg:", []byte(s.Msg), "ip:", s.Ip)
	log.Println("finshed reading msg , read:", s.Msg, "from", s.Ip)
}

// reads non stop reading conncetion from the pkg
// serverNameResponder
func ReadingConnectionContinously(ctx context.Context, conn *net.UDPConn) {
	for {
		conn.SetDeadline(time.Time{})
		log.Println("setting deadline to 0")
		ch := make(chan ReceivedMessage)
		go ReadMsg(conn, ch)
		select {
		case <-ctx.Done():
			log.Println("context canceled", ctx.Err())
			conn.SetDeadline(time.Now())
			log.Println("setting deadline to now")
			return
		case res := <-ch:
			log.Println("returned from channel", res)
			log.Println("ReadingConnectionContinously() compare", strings.Compare(res.Msg, "WhatIsYourName"))
			if res.Msg == "WhatIsYourName" {
				SendMsg(server.Name, res.Ip, ClientPort)
			} else {
				log.Println("received", res.Msg)
			}
		}
	}
}

// Returns IP String
func FindIP(ctx context.Context, conn *net.UDPConn, serverName string) string {
	SendMsg("WhatIsYourName", BroadcastAddr, ServerPort)
	for {
		ch := make(chan ReceivedMessage)
		conn.SetReadDeadline(time.Time{})
		go ReadMsg(conn, ch)
		select {
		case <-ctx.Done():
			log.Println("context canceled", ctx.Err())
			conn.SetDeadline(time.Now())
			return ""
		case res := <-ch:
			log.Println("received FindIP() ", res)
			log.Println("FindIP() received from channel", res)
			log.Println("FindIP() result returned from channel:", res)
			log.Println("serverName:", serverName, "server.Name", server.Name, "res.Msg", res.Msg)
			log.Println("compare", strings.Compare(res.Msg, serverName))
			if res.Msg == serverName {
				log.Println("res.Msg is matching server.Name")
				log.Println("returns", res.Ip)
				close(ch)
				return res.Ip
			}
		}
	}
}
