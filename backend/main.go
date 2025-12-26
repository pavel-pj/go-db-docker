package main

import (
	"database/sql"
	repository "db200/repositories"
	"db200/router"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Подключение к PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Создаем репозиторий
	userRepo := repository.NewUserRepository(db)

	// Создаем роутер
	r := router.NewRouter(userRepo)

	/*
		// Маршруты
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello from Go + PostgreSQL!")
		})

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			if err := db.Ping(); err != nil {
				http.Error(w, "DB connection failed", http.StatusServiceUnavailable)
				return
			}
			fmt.Fprintln(w, "OK")
		})
	*/

	port := "8100"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Server starting on :%s", port)
	// Используй роутер как обработчик!
	log.Fatal(http.ListenAndServe(":"+port, r))
}
