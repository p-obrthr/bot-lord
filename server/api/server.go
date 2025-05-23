package api

import (
	"botlord/bot"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Wrapper struct {
	Bot     *bot.Bot
	Running bool
	mutex   sync.Mutex
}

func NewWrapper() *Wrapper {
	return &Wrapper{}
}

func (wr *Wrapper) StartBot() {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if wr.Running {
		fmt.Println("Bot is already running")
		return
	}
	wr.Bot.Start()
	wr.Running = true
}

func (wr *Wrapper) StopBot() {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if !wr.Running {
		fmt.Println("Bot is not running")
		return
	}
	fmt.Println("Bot stop requested")
	wr.Bot.Stop()
	wr.Running = false
}

func (wr *Wrapper) StartHTTPServer(address string) {
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		wr.StartBot()
		w.Write([]byte("Bot started"))
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		wr.StopBot()
		w.Write([]byte("Bot stopped"))
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if wr.Running {
			w.Write([]byte("running"))
		} else {
			w.Write([]byte("stopped"))
		}
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		maxLogs := 20

		startIndex := len(wr.Bot.Logs) - maxLogs
		if startIndex < 0 {
			startIndex = 0
		}
		visibleLogs := wr.Bot.Logs[startIndex:]

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(visibleLogs)
	})

	fmt.Printf("HTTP server listening on %s\n", address)
	go http.ListenAndServe(address, nil)
}
