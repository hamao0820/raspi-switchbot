package router

import (
	"log"
	"net/http"

	"github.com/hamao0820/raspi-switchbot/switchbot"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// static/index.html を返す
	http.ServeFile(w, r, "static/index.html")
}

func tunrOnHandler(bot *switchbot.SwitchBot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := bot.TurnOn(r.Context()); err != nil {
			log.Printf("failed to turn on the SwitchBot: %v", err)
			http.Error(w, "failed to turn on the SwitchBot", http.StatusInternalServerError)
		}
		w.Write([]byte("SwitchBot turned on\n"))
	}
}
