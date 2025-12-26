/*
package main

import (
	"database/sql"
	repository "db200/repositories"
	"db200/router"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

func main2() {

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

	port := "8100"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Server starting on :%s", port)
	// Используй роутер как обработчик!
	log.Fatal(http.ListenAndServe(":"+port, r))

}
*/
