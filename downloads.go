package server

import (
	"net/http"
	"strings"
)

type downloadHandler struct{}

func (downloadHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(rw, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	if isClientExisting(remoteAddr(r)) == -1 {
		serveConnectHTML(rw, "First Connect", normalConnectNotification())
		return
	}

	if r.URL.Path != "" {
		rw.Header().Set("Content-Disposition:", "attachment; filename=\""+strings.TrimPrefix(r.URL.Path, "/downloads/")+"\"")
	}

	http.FileServer(http.Dir("./uploads")).ServeHTTP(rw, r)
}

func init() {
	router.PathPrefix("/downloads/").Handler(http.StripPrefix("/downloads/", downloadHandler{}))
}
