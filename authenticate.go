package server

import (
	"log"
	"net/http"
)

var (
	err error
)

func init() {
	router.HandleFunc("/authenticate", authenticateHandler)
}

func authenticateHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(rw, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println("Unable To Parse Body : " + r.Method + " " + err.Error())
		http.Error(rw, "UNABLE TO PARSE REQUEST BODY", http.StatusBadRequest)
	}

	log.Println("password = ", r.FormValue("password"))
	if r.FormValue("password") != Password {
		serveConnectHTML(rw, "Authenitcation failed", normalConnectNotification())
		return
	}

	addClient(remoteAddr(r))
	log.Println("Remote Addr Added", remoteAddr(r))
	log.Println("Request ", isClientExisting(remoteAddr(r)))
	listClient()
	serveUploadHandler(rw, ClipboardString)
}
