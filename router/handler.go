package router

import (
	"net/http"

	"github.com/hamao0820/raspi-switchbot/switchbot"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the SwitchBot API"))
}

func tunrOnHandler(bot *switchbot.SwitchBot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := bot.TurnOn(r.Context()); err != nil {
			http.Error(w, "Failed to turn on the SwitchBot: "+err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte("SwitchBot turned on"))
	}
}
