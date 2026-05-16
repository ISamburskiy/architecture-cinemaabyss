package main

import (
	"log"
	"net/http"
)

func main() {
	// Запускаем всех потребителей
	go startMovieConsumer()
	go startUserConsumer()
	go startPaymentConsumer()

	// Настраиваем роутинг
	http.HandleFunc("/api/events/movie", movieHandler)
	http.HandleFunc("/api/events/user", userHandler)
	http.HandleFunc("/api/events/payment", paymentHandler)
	http.HandleFunc("/api/events/health", healthHandler)


	log.Println("Server starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
