package customer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Product представляет товар каталога.
type Customer struct {
	ID        int64
	Email     string
	Nickname  *string
	Age       *int64
	LastLogin *time.Time
	CreatedAt time.Time
}

func AddCustomer(
	ctx context.Context,
	db *sql.DB,
	email string,
	nickname *string,
	age *int64,
	lastLogin *time.Time,
	createdAt time.Time,
) (Customer, error) {

	var c Customer

	// Преобразуем время в строки
	createdAtStr := createdAt.Format(time.RFC3339Nano)

	// Для nullable времени
	var lastLoginStr interface{}
	if lastLogin != nil {
		lastLoginStr = lastLogin.Format(time.RFC3339Nano)
	} else {
		lastLoginStr = nil
	}

	result, err := db.ExecContext(ctx,
		`INSERT INTO customers(email, nickname, age, last_login, created_at) VALUES (?, ?, ?, ?, ?)`,
		email, nickname, age, lastLoginStr, createdAtStr,
	)

	if err != nil {
		return Customer{}, fmt.Errorf("ошибка вставки: %w", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return Customer{}, fmt.Errorf("ошибка получения ID: %w", err)
	}

	// Временные переменные для сканирования
	var (
		nicknameDB     *string
		ageDB          *int64
		lastLoginStrDB *string // сканируем как строку
		createdAtStrDB string  // сканируем как строку
	)

	err = db.QueryRowContext(ctx,
		`SELECT id, email, nickname, age, last_login, created_at FROM customers WHERE id = ?`,
		lastInsertId,
	).Scan(
		&c.ID,
		&c.Email,
		&nicknameDB,     // указатель на string
		&ageDB,          // указатель на int64
		&lastLoginStrDB, // указатель на string (время)
		&createdAtStrDB, // string (время)
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Customer{}, fmt.Errorf("пользователь с ID %d не найден", lastInsertId)
		}
		return Customer{}, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	// Копируем указатели
	c.Nickname = nicknameDB
	c.Age = ageDB

	// Парсим время last_login
	if lastLoginStrDB != nil {
		parsedTime, err := parseTime(*lastLoginStrDB)
		if err == nil {
			c.LastLogin = &parsedTime
		} else {
			// Если не удалось распарсить, используем исходное
			c.LastLogin = lastLogin
		}
	}

	// Парсим created_at
	c.CreatedAt, err = parseTime(createdAtStrDB)
	if err != nil {
		c.CreatedAt = createdAt // fallback
	}

	return c, nil
}

func parseTime(timeStr string) (time.Time, error) {
	// Пробуем RFC3339Nano
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if err == nil {
		return t, nil
	}

	// Пробуем RFC3339
	t, err = time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

	// Пробуем убрать " m=+" часть если есть
	if idx := strings.Index(timeStr, " m=+"); idx != -1 {
		timeStr = timeStr[:idx]
		return time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeStr)
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить время: %s", timeStr)
}

func GetCustomer(ctx context.Context, db *sql.DB, id int64) (Customer, error) {

	var c Customer
	err := db.QueryRowContext(ctx,
		"SELECT id,email,nickname,age, last_login, created_at from customers where id = ?",
		id,
	).Scan(&c.ID, &c.Email, &c.Nickname, &c.Age, &c.LastLogin, &c.CreatedAt)

	if err != nil {
		return Customer{}, err
	}

	return c, nil

}

func ListCustomers(ctx context.Context, db *sql.DB) ([]Customer, error) {

	var customers []Customer

	rows, err := db.QueryContext(ctx,
		"SELECT id,email,nickname,last_login,created_at from customers ORDER BY id",
	)
	if err != nil {
		return customers, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.Email, &c.Nickname, &c.LastLogin, &c.CreatedAt)
		if err != nil {
			return []Customer{}, err
		}
		customers = append(customers, c)

	}

	if err = rows.Err(); err != nil {
		return []Customer{}, err
	}

	return customers, err

}
