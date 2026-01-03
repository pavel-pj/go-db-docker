package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var courses = map[int64]string{
	1: "Introduction to programming",
	2: "Introduction to algorithms",
	3: "Data structures",
}

func main() {

	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, ".log")
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	logger.SetOutput(io.MultiWriter(os.Stdout, file))

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              ":8100",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/courses/description", CourseDescHandler)
	mux.HandleFunc("/sum", SumHandler(logger))

	port := "8100"
	logWithPort := logrus.WithFields(logrus.Fields{
		"port": port,
	})
	logWithPort.Info("Starting a web-server on port")
	logWithPort.Fatal(server.ListenAndServe())

	/*
		port := "8100"
		server := &http.Server{
			Addr:              ":" + port,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Дополнительная информация передается функцией WithFields
		logrus.WithFields(logrus.Fields{
			"port": port,
		}).Info("Starting a web-server on port")
		logrus.Fatal(server.ListenAndServe())
		/*
			err := server.ListenAndServe()
			if err != nil {
				panic(err)
			}

			/*
				port := "8100"

				log.Println("Starting a web-server on port " + port)
				log.Fatal(http.ListenAndServe(":"+port, nil))*/

	//http.ListenAndServe(":8100", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Go to /courses/description"))
	if err != nil {
		log.Printf("welcome to hexlet error: %s\n", err.Error())
	}
}

func CourseDescHandler(w http.ResponseWriter, r *http.Request) {
	getParam := r.URL.Query().Get("course_id")
	param, err := strconv.ParseInt(getParam, 10, 64)
	if err != nil {
		fmt.Fprintf(w, "Parse error: %v", err)
		return
	}
	response, ok := courses[param]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(response))

}

func SumHandler(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		paramX := r.URL.Query().Get("x")
		if paramX == "" {
			http.Error(w, "Missing parameter: x", http.StatusBadRequest)
			return
		}
		paramY := r.URL.Query().Get("y")
		if paramY == "" {
			http.Error(w, "Missing parameter: y", http.StatusBadRequest)
			return
		}
		// Парсим как big.Int (для любых чисел)
		x := new(big.Int)
		_, okX := x.SetString(paramX, 10)
		if !okX {
			http.Error(w, "x should be a valid integer", http.StatusBadRequest)
			return
		}

		y := new(big.Int)
		_, okY := y.SetString(paramY, 10)
		if !okY {
			http.Error(w, "y should be a valid integer", http.StatusBadRequest)
			return
		}

		// Проверяем что числа положительные
		if x.Sign() < 0 || y.Sign() < 0 {
			http.Error(w, "x and y must be positive", http.StatusBadRequest)
			return
		}

		// Складываем
		sum := new(big.Int).Add(x, y)

		// Проверяем не превышает ли MaxInt
		maxInt := big.NewInt(math.MaxInt)
		if sum.Cmp(maxInt) > 0 {
			logger.WithFields(logrus.Fields{
				"x": paramX,
				"y": paramY,
			}).Warn("Sum overflows int")

			// Возвращаем -1
			w.Write([]byte("-1"))
			return
		}

		// Конвертируем big.Int в int (теперь безопасно)
		resultInt := int(sum.Int64())
		result := strconv.Itoa(resultInt)
		w.Write([]byte(result))

	}
}

/*
	dbConn, err := dbInit()
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer dbConn.Close()
	log.Println("✅ Успешное подключение к БД")

	// Создаем экземпляр Queries из sqlc
	productStore := store.NewProductStore(dbConn)
	productService := service.NewProductService(productStore)

	//queries := productsdb.New(dbConn)
	ctx := context.Background()
*/

/*
		// CRUD операции
		created, err := productService.Create(ctx, service.CreateProductInput{
			Slug:        "wooden-desk",
			Title:       "Wooden Desk",
			Description: "Solid oak desk",
			PriceCents:  15000,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created: %v", created)

	p, err := productService.Get(ctx, 26)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(p)

	//err = queries.DeleteAllProducts(ctx)
	/*
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Очистили porducts")
			}


		result, err := queries.CreateProduct(ctx, productsdb.CreateProductParams{
			Slug:        "UUU24",
			Title:       "rqwer",
			Description: "A62562344Q",
			PriceCents:  5342,
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Создание:")
			fmt.Println(result)
		}

		id := result.ID

		result, err = queries.GetProductByID(ctx, id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Show by ID:")
			fmt.Println(result)
		}

		resultIndex, err := queries.ListProducts(ctx, productsdb.ListProductsParams{
			Limit: 10, Offset: 0,
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("INDEX:")
			fmt.Println(resultIndex)
		}

		rowsAffected, err := queries.UpdateProductPrice(ctx, productsdb.UpdateProductPriceParams{
			PriceCents: 999,
			ID:         id,
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("update:")
			fmt.Println(rowsAffected)
		}

		rowsAffected, err = queries.DeleteProduct(ctx, id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("DELETED")
			fmt.Println(rowsAffected)
		}

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


}

/*
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


	return dbConn, nil
}

// Вспомогательная функция для получения переменных окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
*/
