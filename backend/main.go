package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed database/migrations
var migrationsFS embed.FS

func main() {

	dbConn, err := dbInit()
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer dbConn.Close()
	log.Println("✅ Успешное подключение к БД")

	// Создаем экземпляр Queries из sqlc
	//queries := paymentsdb.New(dbConn)
	//ctx := context.Background()

	data, _ := migrationsFS.ReadFile("001_create_users.down.sql")
	fmt.Println(string(data))

	/*
		result, err := queries.CreatePayment(ctx, paymentsdb.CreatePaymentParams{
			InvoiceID:   "inv-42",
			AmountCents: 9900,
			Status:      "pending",
		})
	*/
	/*
		result, err := queries.SetPaymentStatus(ctx, paymentsdb.SetPaymentStatusParams{
			Status: "paid",
			ID:     1,
		})

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}

		/*
			result, err := queries.CreateUser(ctx, userDb.CreateUserParams{
				Email: "Nunuee@mail.ru",
				Name:  "FUFA",
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}

			err = queries.UpdateUserName(ctx, userDb.UpdateUserNameParams{Name: "ЧЕБУРАКА", ID: 1})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Обновлено")
			}

			res, err := queries.DeleteUser(ctx, 4)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}

			/*
				// 2. Создать продукт БЕЗ цены
				_, err = queries.CreateProduct(ctx, productsdb.CreateProductParams{
					Name:   "Компьютер",
					Status: "Active",
					Price:  sql.NullInt32{Int32: 2988, Valid: true},
				})
				if err != nil {
					log.Printf("❌ ОШИБКА: %v", err)
				}

				products, err := queries.GetProducts(ctx)
				if err != nil {
					log.Printf("❌ ОШИБКА: %v", err)
				} else {

					fmt.Println(products)
				}
	*/

	/*
			user, err := queries.CreateUser(ctx, db.CreateUserParams{
				Name:  "Валера Киношников",
				Email: "noneus@mail.ru",
			})

			if err != nil {
				log.Printf("❌ ОШИБКА: %v", err)
			} else {
				log.Printf("✅ Запись юзера: ID=%d", user.ID)
			}


		user, err := queries.GetUserByEmail(ctx, "noneus@mail.ru")
		if err != nil {
			log.Printf("❌ ОШИБКА: %v", err)
		} else {
			log.Println(user)
		}
	*/

}

func dbInit() (*sql.DB, error) {
	// Получаем переменные окружения
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "golang")
	dbPassword := getEnv("DB_PASSWORD", "secret")
	dbName := getEnv("DB_NAME", "app")

	// Формируем connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	//log.Printf("🔗 Connecting to PostgreSQL: %s:%s/%s", dbHost, dbPort, dbName)

	// Подключаемся к БД
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(25)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем соединение с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbConn.PingContext(ctx); err != nil {
		dbConn.Close() // Закрываем при ошибке ping
		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	//log.Println("✅ Connected to PostgreSQL")

	// Запуск миграций
	if err := runMigrations(dbConn); err != nil {
		dbConn.Close() // Закрываем при ошибке миграций
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	return dbConn, nil
}

// Вспомогательная функция для получения переменных окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func runMigrations(db *sql.DB) error {
	// Создаем драйвер для миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	// Создаем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations", // путь к миграциям
		"postgres",                   // имя базы
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	// Запускаем миграции вверх
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations up: %w", err)
	}

	//log.Println("✅ Migrations applied successfully")

	return nil
}

/*
// Обработчик для /api/users (GET и POST)
func usersHandler(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		//case http.MethodGet:
		//handlers.ListUsersHandler(q)(w, r)
		case http.MethodPost:
			handlers.CreateUserHandler(q)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
*/
