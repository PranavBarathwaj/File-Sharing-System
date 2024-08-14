package server

import (
	"io"
	"log"
	"net/http"
	"os"
)

var (
	ClipboardString string
)

func init() {
	router.HandleFunc("/uploadfile", uploadfileHandler)
}

func uploadfileHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(rw, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	if isClientExisting(remoteAddr(r)) == -1 {
		serveConnectHTML(rw, "First Connect", normalConnectNotification())
		return
	}

	err = r.ParseMultipartForm(1 << 30)
	if err != nil {
		http.Error(rw, "UNABLE TO PARSE FORM", http.StatusBadRequest)
		return
	}

	if r.FormValue("clipboardInput") == "" {
		log.Println("Clip Board is empty")
	}
	ClipboardString = r.FormValue("clipboardInput")
	log.Println("Clip Board:", ClipboardString)

	file, handler, err := r.FormFile("file")

	if err != nil && err != http.ErrMissingFile {
		log.Println(err)
		http.Error(rw, "InternalServerError", http.StatusInternalServerError)
		return
	}

	if err == http.ErrMissingFile {
		serveUploadHandler(rw, ClipboardString)
		return
	}

	defer file.Close()

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", 0755)

		if err != nil {
			log.Println("Error in creating Folder " + err.Error())
		}
	}

	localFile, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		log.Println(err)
		http.Error(rw, "InternalServerError", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, file); err != nil {
		log.Println(err)
		http.Error(rw, "InternalServerError", http.StatusInternalServerError)
		return
	}

	serveUploadHandler(rw, ClipboardString)
}
