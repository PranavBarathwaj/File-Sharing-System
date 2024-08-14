package server

import "net/http"

func init() {
	router.HandleFunc("/forgetMe", func(rw http.ResponseWriter, r *http.Request) {
		if isClientExisting(remoteAddr(r)) == -1 {
			serveConnectHTML(rw, "First Connect", normalConnectNotification())
			return
		}

		removeClient(remoteAddr(r))
		rw.WriteHeader(http.StatusAccepted)
		listClient()
		serveConnectHTML(rw, "Thanks", normalConnectNotification())
	})
}
