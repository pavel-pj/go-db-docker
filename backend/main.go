package main

import (
	"context"
	"database/sql"
	c "db200/sql/customer"
	"fmt"
	"log"
	"time"

	//_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//var db *sql.DB

//ЕСЛИ НЕТ СЕРВЕРА ТО ЗАПУСКАТЬ go run main.go

func main() {
	db, err := sql.Open("sqlite", "./test.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 4. ПРЯМОЙ INSERT без всяких функций
	fmt.Println("\n🔄 Пробую DELETE ALL...")
	_, err = db.Exec(
		"DELETE FROM customers",
	)

	if err != nil {
		log.Fatal("❌ DELETE error:", err)
	}
	fmt.Println("\n📋 УДАЛили ВСЕ:")

	// 3. Простой CREATE без IF NOT EXISTS
	_, err = db.Exec(`CREATE TABLE  IF NOT EXISTS  products (
				id INTEGER PRIMARY KEY,
				name TEXT NOT NULL UNIQUE,
				price INTEGER NOT NULL
			)`)
	if err != nil {
		log.Fatal("CREATE products error:", err)
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS users`)
	if err != nil {
		log.Fatal("CREATE users error:", err)
	}

	_, err = db.Exec(`CREATE TABLE  IF NOT EXISTS  users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			status TEXT ,
			age INTEGER  ,
			started_at TIMESTAMPTZ NOT NULL
			)`)
	if err != nil {
		log.Fatal("CREATE users error:", err)
	}

	_, err = db.Exec(`CREATE TABLE  IF NOT EXISTS  customers (
			id INTEGER PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			nickname TEXT,
			age INTEGER,
			last_login TIMESTAMP,
			created_at TIMESTAMP NOT NULL
		)`)
	if err != nil {
		log.Fatal("CREATE customers error:", err)
	}

	ctx := context.Background()
	startedAt := time.Now()
	customer, err := c.AddCustomer(ctx, db, "nome@mail.ru", nil, nil, nil, startedAt)
	if err != nil {
		fmt.Println(err)
	}
	customer, err = c.AddCustomer(ctx, db, "OPPAmail.ru", nil, nil, nil, startedAt)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(customer)
	fmt.Println("Вызов Функции Show: ")
	customer, err = c.GetCustomer(ctx, db, 1)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(customer)

	fmt.Println("Вызов Функции List: ")
	customers, err := c.ListCustomers(ctx, db)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(customers)

	/*
		active := "active"
		startedAt := time.Now()
		_, err = u.AddUser(ctx, db, "Вася", "nome@mail.ru", &active, nil, startedAt)
		if err != nil {
			fmt.Println(err)
		}
		_, err = u.AddUser(ctx, db, "Иван Иваныч", "otto200@mail.ru", &active, nil, startedAt)
		if err != nil {
			fmt.Println(err)
		}
		u.GetAllUsers(ctx, db)
	*/

	/*if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Подключено к SQLite!")
	/*
		// Создание таблицы
		_, err = db.Exec(`
					CREATE TABLE IF NOT EXISTS products(
					id INTEGER PRIMARY KEY,
			    name TEXT NOT NULL UNIQUE,
			    price INTEGER NOT NULL
					)
				`)
		if err != nil {
			log.Fatal(err)
		}

	ctx := context.Background()
	prod, err := p.AddProduct(ctx, db, "AA", 70000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prod)

	/*
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
		}*/

	/*
		rows, err := db.QueryContext(ctx,
			`Select id,name,email from users`,
		)

		if err != nil {
			log.Fatal(err)
		}

		var users []User

		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				log.Fatal(err)
			}

			users = append(users, u)

		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Println(users)
		/*

			email := "john@example.comr"
			u := User{}

			err = db.QueryRowContext(ctx,
				`Select id,name,email from users where email=$1`,
				email,
			).Scan(&u.ID, &u.Name, &u.Email)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("Ошибка: пользователь не найден")
					return
				}
			}

			fmt.Println(u)
	*/
	/*
		res, err := db.ExecContext(ctx,
			`Insert into users (name,email) values ($1,$2)`,
			"Василис", "auto@mail.ru",
		)
		if err != nil {
			log.Fatal(err)
		}

		rows, _ := res.RowsAffected()
		fmt.Println(rows)
	*/
}
