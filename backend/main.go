package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // драйвер PostgreSQL
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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
	orders map[string]Order
}

func (o *OrderStorage) CreateOrder(order Order) (string, error) {
	o.orders[order.ID] = order
	return order.ID, nil

}

// Ошибки
var (
	errOrderNotFound = errors.New("order not found")
)

func (o *OrderStorage) GetOrder(ID string) (Order, error) {

	// Debug - для отладки
	logrus.Debugf("GetOrder called with id: %s", ID)
	logrus.Debugf("Total orders in storage: %d", len(o.orders))

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
