package router

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hamao0820/raspi-switchbot/switchbot"
)

func NewRouter(bot *switchbot.SwitchBot) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", homeHandler)

	return r
}
