package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"tracky/internal/db"
	"tracky/internal/handlers"
	"tracky/internal/tracking"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	connStr := os.Getenv("DB_DSN")
	if connStr == "" {
		log.Fatal("DB_DSN не задан в переменных окружения")
	}
	database, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer database.Close()

	handlers.InitHandlers(database)
	handlers.InitUserHandlers(database)

	// Периодическое обновление статусов
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			tracking.CheckUpdates(database)
			<-ticker.C
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/tracking", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddTrackingHandler(w, r)
		case http.MethodGet:
			handlers.ListTrackingsHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteTrackingHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/tracking/status", handlers.GetTrackingStatusHandler)
	mux.HandleFunc("/api/register", handlers.RegisterUserHandler)
	mux.HandleFunc("/api/carriers", handlers.CarriersHandler)

	log.Println("Tracky server starting on :8080...")
	http.ListenAndServe(":8080", mux)
}
