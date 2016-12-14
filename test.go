package main

import (
	"fmt"
	"proxy/lib/proxy"
	"proxy/lib/sockets"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	sessions := make(sockets.Sessions, 0, 10)
	wg.Add(2)
	proxy := proxy.NewIntercepterProxy("localhost", "8080")
	proxy.Run(wg)

	server, _ := sockets.NewServer("localhost", "2121")
	server.Run(wg)
	for {
		command := <-server.Commands
		switch command {
		case "sess_init\n":
			ses, err := sockets.NewSession(sessions)
                        if err != nil {
                            panic(err)
                        }
			command = fmt.Sprintf("SESSID:%d\nSECRET:%x", ses.ID, ses.Secret)
			fmt.Printf("YAY")
			break
		case "intercept_start":
			proxy.Enabled = true
			break
		case "intercept_stop":
			proxy.Enabled = false
			break
		}
		fmt.Printf("Command: %x - %x", []byte(command), []byte("sess_init\n"))
		fmt.Println(command == "sess_init\n")
		server.Responses <- command

	}
	wg.Wait()
}
