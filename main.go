package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/hamao0820/raspi-switchbot/router"
	"github.com/hamao0820/raspi-switchbot/switchbot"
)

type Config struct {
	Port    string `env:"PORT" envDefault:"8080"`
	Address string `env:"SWITCHBOT_ADDRESS,required"`
}

func main() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := switchbot.ScanSwitchBot(ctx, cfg.Address)
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.NewRouter(bot),
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
