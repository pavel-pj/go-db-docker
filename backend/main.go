package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func main() {
	dsn := "host=localhost port=5450 user=golang password=secret dbname=app sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаём контекст с таймаутом. Если база "зависла", приложение не будет ждать бесконечно.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Именно здесь устанавливается реальное соединение с базой.
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("database unreachable:", err)
	}
	email := "john@example.com"
	u := User{}

	err = db.QueryRowContext(ctx,
		`Select id,name,email from users where email=$1`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email)

	fmt.Println(u)

}
