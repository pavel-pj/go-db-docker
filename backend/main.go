package main

import (
	"context"
	"database/sql"
	p "db200/sql/product"
	"fmt"
	"log"

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
	result, err := db.Exec(
		"DELETE FROM products",
	)

	if err != nil {
		log.Fatal("❌ DELETE error:", err)
	}
	fmt.Println("\n📋 УДАЛили ВСЕ:")

	/*
		// 3. Простой CREATE без IF NOT EXISTS
		_, err = db.Exec(`CREATE TABLE products (
				id INTEGER PRIMARY KEY,
				name TEXT NOT NULL UNIQUE,
				price INTEGER NOT NULL
			)`)

		if err != nil {
			log.Fatal("CREATE error:", err)
		}

		fmt.Println("✅ Таблица создана")
	*/

	// 5. Проверь что в таблице
	rows, _ := db.Query("SELECT * FROM products")
	defer rows.Close()
	fmt.Println("\n📋 Содержимое таблицы:")
	for rows.Next() {
		var id int64
		var name string
		var price int64
		rows.Scan(&id, &name, &price)
		fmt.Printf("  ID: %d, Name: %s, Price: %d\n", id, name, price)
	}

	// 4. ПРЯМОЙ INSERT без всяких функций
	fmt.Println("\n🔄 Пробую INSERT 1...")
	result, err = db.Exec(
		"INSERT INTO products (name, price) VALUES (?, ?)",
		"TEST_122", // гарантированно уникальное
		1000,
	)

	if err != nil {
		log.Fatal("❌ INSERT 1 error:", err)
	}

	id, _ := result.LastInsertId()
	fmt.Printf("✅ INSERT 1 OK, ID: %d\n", id)

	// 6. Попробуй еще один INSERT (должен работать)
	fmt.Println("\n🔄 Пробую INSERT 2...")
	_, err = db.Exec(
		"INSERT INTO products (name, price) VALUES (?, ?)",
		"TEST_215", // другое имя
		2000,
	)

	if err != nil {
		log.Fatal("❌ INSERT 2 error:", err)
	}
	fmt.Println("✅ INSERT 2 OK")

	ctx := context.Background()
	product, err := p.AddProduct(ctx, db, "Валера02", 244)
	if err != nil {
		log.Fatal("❌ INSERT 3 error:", err)
	}
	fmt.Println("✅ INSERT 3 OK")
	fmt.Println(product)

	// 5. Проверь что в таблице
	rows, _ = db.Query("SELECT * FROM products")
	defer rows.Close()
	fmt.Println("\n📋 Содержимое таблицы:")
	for rows.Next() {
		var id int64
		var name string
		var price int64
		rows.Scan(&id, &name, &price)
		fmt.Printf("  ID: %d, Name: %s, Price: %d\n", id, name, price)
	}

	counts, err := p.CountProducts(ctx, db)
	if err != nil {
		log.Fatal("❌ SHOW error:", err)
	}
	fmt.Println("\n📋 Количество записей:")
	fmt.Println(counts)

	products, err := p.ListProducts(ctx, db)
	if err != nil {
		log.Fatal("❌ LIST error:", err)
	}
	fmt.Println("✅ Записи")
	fmt.Println(products)

	/*
		// 6. Попробуй еще один INSERT (должен работать)
		fmt.Println("\n🔄 Пробую INSERT 3..")
		ctx := context.Background()
		_, err = db.ExecContext(ctx,
			"Insert into products (name,price) values(?,?)",
			"ABBB01",
			333444,
		)
		if err != nil {
			log.Fatal("❌ INSERT 3 error:", err)
		}
		fmt.Println("✅ INSERT 3 OK")
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
