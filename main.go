package main

import (
	"github.com/Its-Maniaco/ggcache/cache"
	"github.com/Its-Maniaco/ggcache/server"
)

func main() {
	opts := server.ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	server := server.NewServer(opts, cache.Cacher{})
	server.Start()
}
