package main

import (
	"net/http"
	"fmt"
	"io"
)

func movieHandler(w http.ResponseWriter, r *http.Request) {
	handleRawEvent(w, r, "movie-events")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	handleRawEvent(w, r, "user-events")
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	handleRawEvent(w, r, "payment-events")
}

func handleRawEvent(w http.ResponseWriter, r *http.Request, topic string) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = produceToTopic(topic, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Event sent to %s", topic)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("events-service healthy"))
}