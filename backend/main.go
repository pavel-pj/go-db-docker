package main

import (
	"fmt"
	"time"
"sync"
"net/http"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // драйвер PostgreSQL


	_ "github.com/golang-migrate/migrate/v4/source/file"
)


func fetch (url string, wg *sync.WaitGroup) {


	defer wg.Done()
	 start := time.Now() // Засекаем время начала

	resp,err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка : ",err)
	}

	defer resp.Body.Close()
	fmt.Println(url, resp.Status)
	duration := time.Since(start) // Вычисляем продолжительность
  fmt.Printf("[%v] %s - %s\n", duration, url, resp.Status)

}
 
 

func main() {
var wg sync.WaitGroup

sites:= []string{
	"https://google.com",
	"https://hexlet.io",
	"https://ya.ru",
}

for _,v := range sites {
	fmt.Println("Взяли сайт:%s",v)
	wg.Add(1)
	go func() {
		fetch(v, &wg)
	}()

}
 
	wg.Wait()

	fmt.Println(" Все запросы закончились")
}

/*
type (

	GetTaskResponse struct {
		ID       int64  `json:"id"`
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description" validate:"required,min=3,max=25"`
		Deadline int64  `json:"deadline" validate:"required"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	Task struct {
		ID       int64
		Desc     string
		Deadline int64
	}

)

var (

	taskIDCounter int64 = 1
	tasks               = make(map[int64]Task)

)

var ErrNotFound = fmt.Errorf("Not found model")

//***************************
//JSON

type (

	BinarySearchRequest struct {
		Numbers []int `json:"numbers"`
		Target  int   `json:"target"`
	}

	BinarySearchResponse struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}

)

// ****************************
// Users

	type User struct {
		Email    string
		Name     string
		password string
	}

type (

	UserCreateRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required,min=3,max=50"`
		Password string `json:"password" validate:"required,min=8,max=16"`
	}

	UserCreateReqsponse struct {
		Email string `json:"email" `
		Name  string `json:"name" `
	}

	LoginRequest struct {
		Email    string `json:"email" vaidate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=16"`
	}

	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}

)

var users = map[string]User{}

var (

	errBadCredentials = errors.New("email or password is incorrect")

)

var jwtSignature = []byte("supet-secret-signature-2400")

var contextKeyUser = "user"

// Структура с информацией о фильме

	type Film struct {
		Title    string
		IsViewed bool
	}

// Для простоты описываем хранилище фильмов в коде

	var films = []Film{
		{
			Title:    "The Shawshank Redemption",
			IsViewed: true,
		},
		{
			Title:    "The Godfather",
			IsViewed: true,
		},
		{
			Title:    "The Godfather: Part II",
			IsViewed: false,
		},
	}

type (

	CreateItemRequest struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}

	Item struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}

)

var (

	items []Item

)
*/

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Posts []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title  string
	UserID uint
}
type Product struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string // название товара
	Price      int    // цена
}

type Profile struct {
	ID     int64 // название товара
	UserID int64
	Bio    string
	User   *User `gorm:"foreignKey:UserID;references:ID"` // Обратная связь
}

type Tag struct {
	ID   uint
	Name string
}

type UserPosts struct {
	ID    string
	Name  string
	Count int64
}
 

	//time.Sleep(100 * time.Millisecond)

	/*
			file, err := os.OpenFile(".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				logrus.Fatalf("error opening file: %v", err)
			}
			defer file.Close()

			validate := validator.New()


		// Создание нового логгера с настройками
		newLogger := logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags), // базовый вывод в консоль
			logger.Config{
				SlowThreshold: time.Second, // порог для медленных запросов
				LogLevel:      logger.Info, // подробный уровень логирования
				Colorful:      true,        // цветной вывод для удобства
			},
		)

		// Получаем параметры из переменных окружения
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost" // значение по умолчанию для локальной разработки
		}

		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}

		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "golang"
		}

		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = "secret"
		}

		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "app"
		}

		// Формируем DSN строку
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbName, dbPort,
		)

		logrus.Infof("Подключаемся к БД: %s:%s", dbHost, dbPort)
		// Открытие соединения через драйвер postgres и GORM
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                 newLogger,
			SkipDefaultTransaction: true, // операции без автоматической транзакции
			PrepareStmt:            true, // кэширование подготовленных выражений
		})
		if err != nil {
			// Логирование ошибки подключения и завершение программы
			log.Fatalf("ошибка подключения к базе: %v", err)
		}

		// Если err == nil, соединение успешно установлено
		logrus.Println("Соединение с базой установлено")

		if err := db.AutoMigrate(
			&User{},
			&Product{},
			&Profile{},
			&Tag{},
			&Post{}); err != nil {
			log.Fatalf("ошибка миграции схемы: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("ошибка доступа к пулу соединений: %v", err)
		}

		// Ping проверяет, что соединение живое и база отвечает
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("ошибка пинга базы: %v", err)
		}

		// Максимальное число открытых соединений к базе
		sqlDB.SetMaxOpenConns(10)

		// Максимальное число простаивающих (неиспользуемых) соединений в пуле
		sqlDB.SetMaxIdleConns(5)

		// Максимальное время жизни соединения
		sqlDB.SetConnMaxLifetime(time.Hour)

		/*
			log.Println("Пул соединений настроен и готов к работе")

			tx := db.Session(&gorm.Session{
				DryRun: true, // режим генерации SQL без выполнения
			})

			// Формирование SELECT-запроса без обращения к базе
			stmt := tx.First(&User{}, 1).Statement

			// Вывод текста запроса
			log.Println("Сформированный SQL:", stmt.SQL.String())


		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		tx := db.WithContext(ctx)

		post := Post{
			UserID: 4,
			Title:  "HEllo about World",
		}
		if err = tx.Create(&post).Error; err != nil {
			fmt.Printf("Error:%s", err)
		}

		var uPosts []UserPosts
		err = tx.Raw(`
			select u.id,u."name" , count(p."title") from users u
			left join posts p on u.id = p.user_id
			where p.title like '%Go%'
			group by u.id,u."name"`).Scan(&uPosts).Error

		if err != nil {
			fmt.Printf("Error:%s", err)
		} else {
			fmt.Println(uPosts)
		}

		/*
			err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				user := User{
					Name: "Oleg",
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}

				posts := []Post{
					{Title: "All About Go", UserID: user.ID},
					{Title: "All ToDo Go", UserID: user.ID},
					{Title: "Allilya", UserID: user.ID},
				}
				for _, v := range posts {
					if err = tx.Create(&v).Error; err != nil {
						return err
					}
				}
				return nil
			})
	*/

	/*
		var users1 []User
		if err = db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Order("title DESC").Where("title LIKE ?", "%Go%")
		}).
			Find(&users1).Error; err != nil {
			fmt.Printf("Error : %s", err)
		}

		fmt.Println(users1)
	*/
 
