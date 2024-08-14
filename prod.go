package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("waiting from interupt")
		<-sigChannel
		fmt.Println("Received an interrupt, stopping services...")
		fmt.Println("Shutting Down Server Peacefully")
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Println("Error Occured When Shutting Down the Server : Error : ", err)
		}
		os.Exit(1)
	}()
}
