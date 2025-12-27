package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Product представляет товар каталога.
type User struct {
	ID        int64
	Name      string
	Email     string
	Status    *string
	Age       *string
	StartedAt time.Time
}

// AddProduct создаёт товар и возвращает его с присвоенным идентификатором.
func AddUser(
	ctx context.Context,
	db *sql.DB,
	name string,
	email string,
	status *string,
	age *string,
	startedAt time.Time,
) (User, error) {

	var u User
	startedAtStr := startedAt.Format(time.RFC3339Nano)
	result, err := db.ExecContext(ctx,
		`Insert into users (name,email,status,age,started_at) values(?,?,?,?,?)`,
		name, email, status, age, startedAtStr,
	)

	if err != nil {
		return User{}, err
	}

	lastInsertId, _ := result.LastInsertId()
	var startedAtStrFromDB string
	err = db.QueryRowContext(ctx,
		`Select id,name,email,status,age,started_at from users where id= ?`,
		lastInsertId,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.Age, &startedAtStrFromDB)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("пользователь с ID %d не найден", lastInsertId)
		}
		return User{}, fmt.Errorf("ошибка получения пользователя: %w", err)
	}
	// Теперь парсим строку в time.Time
	u.StartedAt, err = parseTime(startedAtStrFromDB)
	if err != nil {
		// Если не удалось распарсить, используем исходное время
		u.StartedAt = startedAt
	}

	return u, nil
}

func (u User) String() string {
	var statusStr, ageStr string

	if u.Status != nil {
		statusStr = *u.Status
	} else {
		statusStr = "<nil>"
	}

	if u.Age != nil {
		ageStr = *u.Age
	} else {
		ageStr = "<nil>"
	}

	// Форматируем время
	startedAtStr := u.StartedAt.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, Status: %s, Age: %s, StartedAt: %s}",
		u.ID, u.Name, u.Email, statusStr, ageStr, startedAtStr)
}

func GetAllUsers(ctx context.Context, db *sql.DB) error {
	rows, err := db.Query("SELECT id, name, email, status, age, started_at FROM users")
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()

	fmt.Println("\n📋 Содержимое таблицы users:")
	for rows.Next() {
		var u User
		var startedAtStr string

		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.Age, &startedAtStr)
		if err != nil {
			return fmt.Errorf("ошибка сканирования: %w", err)
		}

		// Выводим сырую строку для отладки
		//fmt.Printf("DEBUG: started_at из БД: %q\n", startedAtStr)

		// Пробуем разные форматы
		u.StartedAt, err = parseTime(startedAtStr)
		if err != nil {
			fmt.Printf("⚠️ Ошибка парсинга времени: %v\n", err)
		}

		fmt.Println(u.String())
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// Функция для парсинга времени из SQLite
func parseTime(timeStr string) (time.Time, error) {
	// Пробуем разные форматы
	formats := []string{
		"2006-01-02 15:04:05.999999999 -0700 MST m=+0.000000001", // Ваш формат
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999999 -0700",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05Z",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить время: %s", timeStr)
}

func ListUsersByStatuses(ctx context.Context, db *sql.DB, status []string, order string) ([]User, error) {

	var users []User

	allowed := map[string]string{
		"id_asc":    "id ASC",
		"name_asc":  "name ASC",
		"name_desc": "name DESC",
	}

	params := strings.TrimRight(strings.Repeat("?,", len(status)), ",")

	args := make([]interface{}, len(status))
	for i, v := range status {
		args[i] = v
	}

	orderBy, exists := allowed[order]
	if !exists {
		orderBy = allowed["id_asc"]
	}

	if len(status) == 0 {
		return users, nil
	}

	query := fmt.Sprintf("SELECT id,name,email,status from users where status IN (%s) ORDER BY %s", params, orderBy)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return users, nil
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Status)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil

}
