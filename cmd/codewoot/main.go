package main

import (
	"log"

	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/server"
)

func main() {
	cfg, _ := config.Load()
	srv := server.New(cfg)
	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
