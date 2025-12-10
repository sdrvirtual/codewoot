package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/server"
)

func main() {
	ctx := context.Background()
	cfg, _ := config.Load()
	dbCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		log.Fatalln(err)
	}

	dbCfg.MaxConns = 10
	dbCfg.MinConns = 2
	dbCfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.Close()

	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		log.Fatalln(err)
	}
	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		log.Fatalln(err)
	}
	db.Close()

	srv := server.New(cfg, pool)
	defer srv.Close()

	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
