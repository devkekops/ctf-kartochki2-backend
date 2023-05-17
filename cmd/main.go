package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env"

	"github.com/devkekops/ctf-kartochki2-backend/internal/app/config"
	"github.com/devkekops/ctf-kartochki2-backend/internal/app/server"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.SecretKey, "k", cfg.SecretKey, "secret key")
	flag.StringVar(&cfg.LicenseKeyHash, "h", cfg.LicenseKeyHash, "license key hash")
	flag.Parse()

	log.Fatal(server.Serve(&cfg))
}
