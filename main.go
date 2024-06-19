package main

import (
	"log"
	"net"
	"time"

	"github.com/Its-Maniaco/ggcache/cache"
	"github.com/Its-Maniaco/ggcache/server"
)

func main() {
	opts := server.ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", opts.ListenAddr)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET Foo bar 300000000000"))

		time.Sleep(time.Second * 2)

		conn.Write([]byte("GET Foo"))
	}()

	server := server.NewServer(opts, cache.New())
	server.Start()
}
