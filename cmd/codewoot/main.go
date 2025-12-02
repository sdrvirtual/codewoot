package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/server"
	"github.com/sdrvirtual/codewoot/internal/services"
)

func Cu() {
	cfg, _ := config.Load()
	service := services.NewChatwootService(cfg)

	contact := services.ContactInfo{
		Name: "Cássio Ávila",
		Phone: "+5534988781744",
	}

	ctx := context.TODO()

	ctt, err := service.SetupContact(ctx, &contact)
	if err != nil {
		log.Fatal(err)
	}
	s, _ := json.MarshalIndent(ctt, "", "\t")
	fmt.Println(string(s))
}

func main() {
	cfg, _ := config.Load()
	srv := server.New(cfg)
	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
