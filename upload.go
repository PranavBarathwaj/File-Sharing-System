package server

import (
	"html/template"
	"net/http"

	"github.com/File-Sharing-System/static"
)

var (
	uploadTemplate *template.Template
)

func init() {
	router.HandleFunc("/upload", uploadHandler)
	uploadTemplate, err = template.New("uploadtemplate").Parse(static.UploadHTML)
	if err != nil {
		panic("Error Parsing Connect HTML" + err.Error())
	}
}

func uploadHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(rw, "BAD REQUEST", http.StatusBadRequest)
		return
	}
	if isClientExisting(remoteAddr(r)) == -1 {
		serveConnectHTML(rw, "First Connect", normalConnectNotification())
		return
	}
	serveUploadHandler(rw, ClipboardString)
}

func serveUploadHandler(rw http.ResponseWriter, clipboard string) {
	uploadTemplate.Execute(rw, struct {
		Clipboard string
	}{
		clipboard,
	})
}
