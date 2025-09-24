package router

import (
	"net/http"
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
	r.Post("/api/turn_on", tunrOnHandler(bot))

	// 静的ファイルを配信
	fileServer := http.FileServer(http.Dir("static/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return r
}
