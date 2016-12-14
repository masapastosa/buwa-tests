package sockets

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// intercept_start : Start intercepter proxy
// intercept_stop  : Stop intercepter proxy

type Server struct {
	Addr      string
	Port      string
	Commands  chan string
	Responses chan string
	Listener  net.Listener
}

func (s *Server) Run(wg sync.WaitGroup) error {
	ln, err := net.Listen("tcp", s.Addr+":"+s.Port)
	// Listener.Close ?
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				// TODO: Handle error
				fmt.Println("Error")
			}
			go s.handleConn(conn)
		}
		wg.Done()
	}()
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		// TODO: Mirar si se leen > 1024 bytes
		_, err := conn.Read(buf)
		if err != nil {
			return
		}
		// Send command and get the response, both read and write from
		// chan is blocking

                // Ugly but it just wroks
                length := 0
                for i := 0; i < len(buf) && buf[i] != '\x00'; i++ {
                    length++
                }
                buf2 := make([]byte, length)
                for i := 0; i < length; i++ {
                    buf2[i] = buf[i]
                }
                // ----

		s.Commands <- string(buf2)
		response := <-s.Responses

		_, err = conn.Write([]byte(response))
		if err != nil {
			return
		}
	}
}

func NewServer(addr string, port string) (*Server, error) {
	ret := &Server{}
	if port == "" {
		return nil, errors.New("server: Port must be specified")
	}
	ret.Addr = addr
	ret.Port = port
	// Revisar esto
	ret.Commands = make(chan string)
	ret.Responses = make(chan string)

	return ret, nil
}
