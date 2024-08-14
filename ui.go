package ui

import (
	"context"
	"log"
	"time"

	"github.com/File-Sharing-System/server"
	"github.com/File-Sharing-System/static"
	"github.com/File-Sharing-System/udp"
	"github.com/zserge/lorca"
)

func StartUI() {
	if lorca.ChromeExecutable() == "" {
		lorca.PromptDownload()
		return
	}

	ui, _ := lorca.New("data:text/html,"+static.UiIndexHTML, "", 480, 320)
	defer ui.Close()

	ui.Bind("main", func() {
		ui.Load("data:text/html," + static.UiIndexHTML)
	})

	ui.Bind("createServerHTML", func() {
		ui.Load("data:text/html," + static.UiServerHTML)
	})

	ui.Bind("connectToServerHTML", func() {
		ui.Load("data:text/html," + static.UiClientHTML)
	})

	ui.Bind("createServer", func(sn, sp string) {
		server.Name = sn
		server.Password = sp
		var ip string
		udp.CreateServerConn()
		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(time.Duration(udp.WaitingTime)*time.Second, cancel)
		ip = udp.FindIP(ctx, udp.ServerConn, sn)
		if ip != "" {
			return
		}
		go func() {
			udp.ReadingConnectionContinously(context.Background(), udp.ServerConn)
		}()
		go func() {
			server.StartServer()
		}()
		ui.Load("http://localhost:8080/connect")
	})

	ui.Bind("clientConnect", func(sn string) {
		if sn == "" {
			log.Println("sn is \"\"")
			return
		}
		var ip string

		udp.CreateClientConn()
		ip = udp.FindIP(context.Background(), udp.ClientConn, sn)
		if ip == "" {
			log.Println("ip not found", ip)
			return
		}
		ui.Load("http://" + ip + ":8080/connect")
	})
	defer func() {
		if server.ServerOn {
			server.Shutdown()
			udp.CloseClientConn()
			udp.CloseSeverConn()
		}
	}()

	<-ui.Done()
}
