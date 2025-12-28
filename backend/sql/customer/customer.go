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

type UserUpdate struct {
	Email string
	Age   int64
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
		"SELECT id, email, nickname, age, last_login, created_at FROM customers ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Customer

		// Просто сканируем прямо в указатели
		// Драйвер сам установит nil для NULL значений
		err := rows.Scan(
			&c.ID,
			&c.Email,
			&c.Nickname,  // *string
			&c.Age,       // *int64
			&c.LastLogin, // *time.Time
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func LoopPrepared(ctx context.Context, db *sql.DB, toUpdate []UserUpdate) error {
	stmt, err := db.PrepareContext(ctx,
		"UPDATE customers set age= ? where email =?",
	)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, v := range toUpdate {
		_, err := stmt.ExecContext(ctx, v.Age, v.Email)
		if err != nil {
			return err
		}
	}
	return nil

}

func (c Customer) String() string {
	// Преобразуем указатели в строки
	ageStr := "<nil>"
	if c.Age != nil {
		ageStr = fmt.Sprintf("%d", *c.Age)
	}

	nicknameStr := "<nil>"
	if c.Nickname != nil {
		nicknameStr = *c.Nickname
	}

	lastLoginStr := "<nil>"
	if c.LastLogin != nil {
		lastLoginStr = c.LastLogin.Format(time.RFC3339)
	}

	return fmt.Sprintf("Customer{ID:%d Email:%s Age:%s Nickname:%s LastLogin:%s CreatedAt:%v}",
		c.ID, c.Email, ageStr, nicknameStr, lastLoginStr, c.CreatedAt)
}

func LoopShow(ctx context.Context, db *sql.DB, ages []int64) error {

	stmt, err := db.PrepareContext(ctx,
		"SELECT id,email,age from customers where age > ?",
	)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, a := range ages {

		rows, err := stmt.QueryContext(ctx, a)
		if err != nil {
			return err
		}
		fmt.Printf("Выборка для age=%d: \n", a)
		for rows.Next() {
			var (
				id    int64
				email string
				age   *int64
			)

			rows.Scan(&id, &email, &age)
			ageN := "<nil>"
			if age != nil {
				ageN = fmt.Sprintf("%d", *age)
			}

			fmt.Printf("id: %d, email:%s, age: %s\n", id, email, ageN)
		}
		rows.Close()

		if err = rows.Err(); err != nil {
			return err
		}

	}

	return nil

}
