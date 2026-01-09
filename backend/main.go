package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	jwt "github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq" // драйвер PostgreSQL
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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

func main() {

	file, err := os.OpenFile(".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	validate := validator.New()
	/*
		vErr := validate.RegisterValidation("allowable_country", func(fl validator.FieldLevel) bool {
			// Проверяем страну
			text := fl.Field().String()
			for _, country := range ValidCountry {
				if country == text {
					return true
				}
			}
			return false
		})
		if vErr != nil {
			log.Fatal("register validation ", vErr)
		}

		taskHandler := &TaskHandler{
			storage: &TaskStorage{
				tasks: tasks,
			},
			validator: validate,
		}
	*/
	authHandler := &AuthHandler{
		storage: &AuthStorage{
			users: users,
		},
		validator: validate,
	}
	userHandler := &UserHandler{
		storage: &AuthStorage{
			users: users,
		},
	}

	webApp := fiber.New()

	publicGroup := webApp.Group("")
	publicGroup.Post("/register", authHandler.CreateUser)
	publicGroup.Post("/login", authHandler.Login)

	authorizedGroup := webApp.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: jwtSignature,
		},
		ContextKey: contextKeyUser,
	}))

	authorizedGroup.Get("/profile", userHandler.Profile)

	/*
		webApp.Use(limiter.New(limiter.Config{
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			Max:        1,
			Expiration: 2 * time.Second,
		}))

		webApp.Use(requestid.New())
		webApp.Use(logger.New(logger.Config{
			Format: "${locals:requestid}: ${method} ${path} - ${status} \n",
			Output: file,
		}))
	*/

	//webApp.Post("/tasks", authHandler.CreateUser)

	//webApp.Post("/tasks", taskHandler.CreateTask)
	//webApp.Get("/tasks/:id", taskHandler.GetTask)
	//webApp.Patch("/tasks/:id", taskHandler.UpdateTask)
	//webApp.Delete("/tasks/:id", taskHandler.DeleteTask)

	logrus.Fatal(webApp.Listen(":8100"))

}

type AuthHandler struct {
	storage   *AuthStorage
	validator *validator.Validate
}

type UserHandler struct {
	storage *AuthStorage
}

type AuthStorage struct {
	users map[string]User
}

func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	var request UserCreateRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Error: %w", err))
	}

	err = h.validator.Struct(request)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error:":   "bad validateion",
			"details:": err.Error(),
		})
	}

	_, exists := h.storage.users[request.Email]
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error:": "user is already registered",
		})
	}

	h.storage.users[request.Email] = User{
		Email:    request.Email,
		Name:     request.Name,
		password: request.Password,
	}

	return c.Status(fiber.StatusOK).JSON(UserCreateReqsponse{
		Name:  request.Name,
		Email: request.Email,
	})

}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request LoginRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.validator.Struct(request)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message:": "Error",
			"details":  err.Error(),
		})
	}

	user, exists := h.storage.users[request.Email]
	if !exists {
		return errBadCredentials
	}

	if user.password != request.Password {
		return errBadCredentials
	}

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSignature)
	if err != nil {
		return c.SendStatus(fiber.StatusOK)
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{
		AccessToken: t,
	})

}

func getUserPayloadFromCtx(c *fiber.Ctx) (jwt.MapClaims, bool) {
	jwtToken, ok := c.Context().Value(contextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": c.Context().Value(contextKeyUser),
		}).Error("wrong type of JWT token in context")
		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")
		return nil, false
	}

	return payload, true

}

// Структура HTTP-ответа с информацией о пользователе
type ProfileResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (u *UserHandler) Profile(c *fiber.Ctx) error {
	payload, ok := getUserPayloadFromCtx(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userInfo, ok := u.storage.users[payload["sub"].(string)]
	if !ok {
		return errors.New("user not found")
	}

	return c.Status(fiber.StatusOK).JSON(ProfileResponse{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	})
}

/*
	webApp.Post("/search", func(c *fiber.Ctx) error {
		var request BinarySearchRequest
		err := c.BodyParser(&request)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(BinarySearchResponse{
				TargetIndex: targetNotFound,
				Error:       "Invalid JSON",
			})
		}

		targetIndex := slices.Index(request.Numbers, request.Target)
		if targetIndex == -1 {
			return c.Status(fiber.StatusNotFound).JSON(BinarySearchResponse{
				TargetIndex: targetNotFound,
				Error:       "Invalid JSON",
			})
		}

		return c.Status(fiber.StatusOK).JSON(BinarySearchResponse{
			TargetIndex: targetIndex,
		})

	})

	// Оборачиваем в функцию логирования, чтобы видеть ошибки, если они возникнут
	logrus.Fatal(webApp.Listen(":8100"))

}

type TaskStorageInterface interface {
	CreateTask(task Task) int64
	GetTask(id int64) (Task, error)
	UpdateTask(task Task) (Task, error)
	DeleteTask(id int64) error
}

type TaskHandler struct {
	storage   TaskStorageInterface
	validator *validator.Validate
}

type TaskStorage struct {
	tasks map[int64]Task
}

func (t *TaskHandler) CreateTask(c *fiber.Ctx) error {

	var request CreateTaskRequest
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	err = t.validator.Struct(request)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	task := Task{
		ID:       taskIDCounter,
		Desc:     request.Desc,
		Deadline: request.Deadline,
	}

	id := t.storage.CreateTask(task)
	taskIDCounter++

	return c.Status(201).JSON(CreateTaskResponse{
		ID: id,
	})

}

func (t *TaskHandler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	task, err := t.storage.GetTask(idInt64)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(200).JSON(GetTaskResponse{
		ID:       task.ID,
		Desc:     task.Desc,
		Deadline: task.Deadline,
	})

}

func (t *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Id should exists")
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error Json")
	}

	var request UpdateTaskRequest
	err = c.BodyParser(&request)
	if err != nil {
		return c.Status(400).SendString("Error Json")
	}

	task := Task{
		ID:       idInt64,
		Desc:     request.Desc,
		Deadline: request.Deadline,
	}

	_, err = t.storage.UpdateTask(task)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)

}

func (t *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Id should exists")
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error ID")
	}

	err = t.storage.DeleteTask(idInt64)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)

}

func (t *TaskStorage) CreateTask(task Task) int64 {
	t.tasks[task.ID] = task
	return task.ID

}

func (t *TaskStorage) GetTask(id int64) (Task, error) {

	result, ok := t.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return result, nil

}

func (t *TaskStorage) UpdateTask(task Task) (Task, error) {
	t.tasks[task.ID] = task
	updated, ok := t.tasks[task.ID]
	if !ok {
		return Task{}, fmt.Errorf("Item not found")
	}

	return updated, nil

}

func (t *TaskStorage) DeleteTask(id int64) error {

	_, ok := t.tasks[id]
	if !ok {
		return fmt.Errorf("Item not found")
	}
	delete(t.tasks, id)

	return nil

}

/*
type Link struct {
	External string
	Internal string
}

type LinkStorageInterface interface {
	CreateLink(link Link) error
	GetLink(external string) (Link, error)
}

type LinkHandler struct {
	storage LinkStorageInterface
}

func (l *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest
	err := c.BodyParser(&request)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	if request.External == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON\n")
	}

	if request.Internal == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON\n")
	}
	link := Link{
		External: request.External,
		Internal: request.Internal,
	}
	err = l.storage.CreateLink(link)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Error: %v", err))
	}

	return c.Status(200).SendString("Отлично")

}

func (l *LinkHandler) GetLink(c *fiber.Ctx) error {

	external := c.Params("external")
	if external == "" {
		return c.Status(fiber.StatusBadRequest).SendString("External parameter is required")
	}

	// ДЕКОДИРУЕМ его
	decoded, err := url.QueryUnescape(external)
	if err != nil {
		return c.Status(400).SendString("Неправильный URL")
	}

	link, err := l.storage.GetLink(decoded)
	if err != nil {
		if errors.Is(err, ErrorNotFound) {
			return c.Status(fiber.StatusNotFound).SendString("Link not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return c.Status(200).JSON(GetLinkResponse{
		Internal: link.Internal,
	})

}

type LinkStorage struct {
	links map[string]string
}

func (s *LinkStorage) CreateLink(link Link) error {
	s.links[link.External] = link.Internal
	return nil

}

var (
	ErrorNotFound = errors.New("Link not found")
)

func (s *LinkStorage) GetLink(external string) (Link, error) {
	result, ok := s.links[external]
	if !ok {
		return Link{}, ErrorNotFound
	}

	link := Link{
		External: external,
		Internal: result,
	}

	return link, nil
}

/*
type (
	CreateOrderRequest struct {
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}

	CreateOrderResponse struct {
		ID string `json: "id"`
	}

	GetOrderResponse struct {
		ID         string  `json:"id"`
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}
)

func main() {

	// Устанавливаем вывод в stdout и формат
	logrus.SetOutput(os.Stdout) // ← ВАЖНО!
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true, // Цвета для читаемости
	})

	logrus.Info("=== Server starting ===")

	OrderHandler := &OrderHandler{
		storage: &OrderStorage{
			orders: make(map[string]Order),
		},
	}

	webApp := fiber.New()
	webApp.Post("orders", OrderHandler.CreateOrder)
	webApp.Get("orders/:id", OrderHandler.GetOrder)

	// Оборачиваем в функцию логирования, чтобы видеть ошибки, если они возникнут
	logrus.Fatal(webApp.Listen(":8100"))
}

type OrderCreatorGetter interface {
	CreateOrder(order Order) (string, error)
	GetOrder(ID string) (Order, error)
}

type OrderHandler struct {
	storage OrderCreatorGetter
}

func (o *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return fmt.Errorf("Erorr: %w", err)
	}

	order := Order{
		ID:         uuid.New().String(),
		UserID:     request.UserID,
		ProductIDs: request.ProductIDs,
	}

	logrus.WithFields(logrus.Fields{
		"user_id":     request.UserID,
		"product_ids": request.ProductIDs,
	}).Debug("Parsed request")

	ID, err := o.storage.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("Erorr: %w", err)
	}

	return c.JSON(CreateOrderResponse{
		ID: ID,
	})

}

func (o *OrderHandler) GetOrder(c *fiber.Ctx) error {

	Id := c.Params("id")
	if Id == "" {
		return fmt.Errorf("Требуется ввести id")
	}
	fmt.Printf("[DEBUG] CreateOrder")
	logrus.Info("=== Server starting ===")

	order, err := o.storage.GetOrder(Id)
	if err != nil {
		return err
	}

	return c.JSON(GetOrderResponse(order))

}

type Order struct {
	ID         string
	UserID     int64
	ProductIDs []int64
}

type OrderStorage struct {
	mu     sync.Mutex
	orders map[string]Order
}

func (o *OrderStorage) CreateOrder(order Order) (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.orders[order.ID] = order
	return order.ID, nil

}

// Ошибки
var (
	errOrderNotFound = errors.New("order not found")
)

func (o *OrderStorage) GetOrder(ID string) (Order, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	order, exists := o.orders[ID]
	if !exists {
		return Order{}, errOrderNotFound
	}
	return order, nil

}

// Вспомогательная функция для получения ключей
func getMapKeys(m map[string]Order) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

*/
