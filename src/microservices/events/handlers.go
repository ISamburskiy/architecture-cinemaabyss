package main

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status": "success",
		"message": fmt.Sprintf("Event sent to %s", topic),
		"topic":   topic,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}