package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	port string = "8080"
)

var (
	server   *http.Server
	Name     string = ""
	Password string = ""
	ServerOn bool   = false
)

func init() {
	server = &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Handler:      router,
	}
}

func Shutdown() {
	log.Println("Shutting Down Server")
	server.Shutdown(context.Background())
}

func StartServer() {
	ServerOn = true
	if Name == "" || Password == "" {
		panic("no server name or password")
	}
	log.Println("Starting Server with name", Name, "password", Password)
	server.ListenAndServe()
}
