package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Its-Maniaco/ggcache/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, cr cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      cr,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("server listen error: %s", err)
	}
	log.Printf("server starting on port: %s", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("listener accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}

}

func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))
	}
}
