package server

import (
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/File-Sharing-System/static"
)

var (
	connectTemplate *template.Template
	ip              string
)

func init() {
	router.HandleFunc("/connect", connectHandler)
	connectTemplate, err = template.New("connecttemplate").Parse(static.ConnectHTML)
	if err != nil {
		panic("Error Parsing Connect HTML" + err.Error())
	}
	ip = getOutboundIP().String()
}

func connectHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(rw, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	serveConnectHTML(rw, "Connect Here", normalConnectNotification())

}

func serveConnectHTML(rw http.ResponseWriter, msg, notif string) {
	connectTemplate.Execute(rw, struct {
		Message      string
		Notification string
	}{
		msg,
		notif,
	})
}

func normalConnectNotification() string {
	return "For getting connect to this server with other devices visit http://" + ip + ":8080/connect in your device browser"
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
