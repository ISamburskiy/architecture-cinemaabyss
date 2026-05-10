package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"io"
)


type Proxy struct {
	monolithURL      string
	moviesServiceURL string
	migrationPercent int
	eventsServiceURL string
}

func NewProxy() *Proxy {
    percent := 50 // по умолчанию

    if envPercent, exists := os.LookupEnv("MOVIES_MIGRATION_PERCENT"); exists {
        if p, err := strconv.Atoi(envPercent); err == nil && p >= 0 && p <= 100 {
            percent = p
        }
    }

    return &Proxy{
        monolithURL:      getEnv("MONOLITH_URL", "http://localhost:8080"),
        moviesServiceURL: getEnv("MOVIES_SERVICE_URL", "http://localhost:8081"),
        migrationPercent:   percent,
		eventsServiceURL: getEnv("EVENTS_SERVICE_URL", "http://localhost:8082"),
    }
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (p *Proxy) shouldRouteToMoviesService() bool {
	return rand.Intn(100) < p.migrationPercent
}

func (p *Proxy) handleMovies(w http.ResponseWriter, r *http.Request) {
	targetURL := p.monolithURL
	if p.shouldRouteToMoviesService() {
		targetURL = p.moviesServiceURL
	}

	proxyReq, err := http.NewRequest(r.Method, targetURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		fmt.Printf("Error copying response: %v\n", err)
	}
}

func (p *Proxy) handleMonolithRequests(w http.ResponseWriter, r *http.Request) {
	targetURL := p.monolithURL

	proxyReq, err := http.NewRequest(r.Method, targetURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		fmt.Printf("Error copying response: %v\n", err)
	}
}

func (p *Proxy) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"migration_percent": p.migrationPercent,
		"routes": map[string]string{
			"monolith": p.monolithURL,
		"movies_service": p.moviesServiceURL,
	},
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	proxy := NewProxy()

	//health
	http.HandleFunc("/health", proxy.healthHandler)
	// migrated movies endpoints
	http.HandleFunc("/api/movies", proxy.handleMovies)
	http.HandleFunc("/api/movies/", proxy.handleMovies)
	// monolith only endpoints
	http.HandleFunc("/api/users", proxy.handleMonolithRequests)
	http.HandleFunc("/api/payments", proxy.handleMonolithRequests)
	http.HandleFunc("/api/subscriptions", proxy.handleMonolithRequests)
	// events tbd?


	port := getEnv("PORT", "8000")
	fmt.Printf("Proxy server starting on port %s (migration: %d%%)\n", port, proxy.migrationPercent)
	fmt.Println("Monolith:", proxy.monolithURL)
	fmt.Println("Movies Service:", proxy.moviesServiceURL)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
