package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hamao0820/raspi-switchbot/router"
)

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.NewRouter(),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Printf("Starting server on port :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
