package server

import (
	"log"
	"net/http"
	"strings"
)

var (
	clientList []string
)

func remoteAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func addClient(remoteAddr string) {
	if isClientExisting(remoteAddr) == -1 {
		clientList = append(clientList, remoteAddr)
	}
}

func isClientExisting(remoteAddr string) int {
	for i, c := range clientList {
		if c == remoteAddr {
			return i
		}
	}
	return -1
}

func listClient() {
	log.Println("--- Client List ---")
	for i, c := range clientList {
		log.Println("Client", i+1, ":", c)
	}
	log.Println("--- List End ---")
}

func removeClient(remoteAddr string) {
	pos := isClientExisting(remoteAddr)
	if pos != -1 {
		clientList[pos] = clientList[len(clientList)-1]
		clientList = clientList[:len(clientList)-1]
	}
}
