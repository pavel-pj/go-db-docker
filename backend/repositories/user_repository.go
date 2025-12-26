package repository

import (
	"database/sql"
	"db200/models"
	"fmt"
	"html"
	"strings"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(user *models.User) (int, error) {
	// Очищаем входные данные от потенциального XSS
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))

	query := `
        INSERT INTO users (name, email) 
        VALUES ($1, $2) 
        RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(query, user.Name, user.Email).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}
