package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Its-Maniaco/ggcache/cache"
	cmnd "github.com/Its-Maniaco/ggcache/cmd"
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
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err)
			break
		}

		msg := buf[:n]
		fmt.Println(string(msg))
		go s.handleCmd(conn, buf[:n])
	}
}

func (s *Server) handleCmd(conn net.Conn, rawCMD []byte) {
	msg, err := cmnd.ParseMessage(rawCMD)
	if err != nil {
		fmt.Println("failed to Parse command: ", err)
		conn.Write([]byte(err.Error()))
		return
	}

	switch msg.Cmd {
	case cmnd.CMDSet:
		err = s.handleSetCmd(conn, msg)
	case cmnd.CMDGet:
		err = s.handleGetCmd(conn, msg)
	}
	// if switch returns an err
	if err != nil {
		fmt.Println("failed to Handle command: ", err)
		conn.Write([]byte(err.Error()))
	}
}

func (s *Server) handleSetCmd(conn net.Conn, msg *cmnd.Message) error {
	if err := s.cache.Set(msg.Key, msg.Val, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) handleGetCmd(conn net.Conn, msg *cmnd.Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		fmt.Println("HERE????")
		return err
	}
	fmt.Println("OR HERE????")
	_, err = conn.Write(val)
	return err
}

// get the command and distribute it to followers
// ctx to set deadline, sync b/w followers if deleted or not
// if follower is out of sync, tell it the message it should delete
func (s *Server) sendToFollowers(ctx context.Context, msg *cmnd.Message) error {
	return nil
}
