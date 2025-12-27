package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Product представляет товар каталога.
type Product struct {
	ID    int64
	Name  string
	Price int64
}

// AddProduct создаёт товар и возвращает его с присвоенным идентификатором.
func AddProduct(ctx context.Context, db *sql.DB, name string, price int64) (Product, error) {

	var p Product
	result, err := db.ExecContext(ctx,
		`Insert into products (name,price) values(?,?)`,
		name, price,
	)

	if err != nil {
		return Product{}, err
	}

	lastInsertId, _ := result.LastInsertId()
	fmt.Println(lastInsertId)
	err = db.QueryRowContext(ctx,
		`Select id,name,price from products where id= ?`,
		lastInsertId,
	).Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Ошибка. Не существует записи с таким ID")
		}
	}

	return p, nil
}

func CountProducts(ctx context.Context, db *sql.DB) (int64, error) {
	var count int64
	err := db.QueryRowContext(ctx,
		"SELECT count(id) from products",
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil

}

func ListProducts(ctx context.Context, db *sql.DB) ([]Product, error) {
	var products []Product
	rows, err := db.QueryContext(ctx,
		"Select id,name,price from products order by id DESC",
	)
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return products, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}
