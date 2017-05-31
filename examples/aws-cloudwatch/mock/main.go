package main

import (
	"github.com/anarcher/mockingjay/pkg/log"
	"github.com/asdine/storm"

	"net/http"
	"os"
)

type Handler struct {
	mux map[string]func(http.ResponseWriter, *http.Request)
	db  *storm.DB
}

func main() {
	logger := log.Logger

	db, err := storm.Open("metrics.db")
	defer db.Close()

	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	handler := &Handler{
		mux: make(map[string]func(http.ResponseWriter, *http.Request)),
		db:  db,
	}
	handler.mux["GetMetricStatistics"] = handler.GetMetricStatistics
	handler.mux["SetDesiredCapacity"] = handler.SetDesiredCapacity

	server := http.Server{
		Addr:    ":8081",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Log("err", err)
	}

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger
	if err := r.ParseForm(); err != nil {
		logger.Log("err", err)
		http.Error(w, "ParseForm Error", 500)
	}
	action := r.FormValue("Action")
	logger.Log("action", action)
	if h, ok := h.mux[action]; ok {
		h(w, r)
		return
	} else {
		http.Error(w, "Action not found", 500)
	}
}
