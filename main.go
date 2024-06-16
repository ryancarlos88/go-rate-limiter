package main

import (
	"github.com/ryancarlos88/go-rate-limiter/config"
	"github.com/ryancarlos88/go-rate-limiter/internal/web/server"
)

func main() {

	cfg := config.NewConfig()

	s := server.NewServer(cfg)

	s.Start()
}
