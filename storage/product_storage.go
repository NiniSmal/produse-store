package storage

import (
	"database/sql"
	"time"
)

type Product struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	UnitPrice    float64   `json:"unit_price"`
	UnitsInStock int64     `json:"units_in_stock"`
	CreatedAt    time.Time `json:""`
	UpdatedAt    time.Time `json:"updated_at"`
	CategoryID   int64     `json:"category_id"`
}
type ProductStorage struct {
	db *sql.DB
}

func NewProductStorage(p *sql.DB) *ProductStorage {
	return &ProductStorage{
		db: p,
	}
}

func (p *ProductStorage) Save(product Product) (int64, error) {
	query := "INSERT INTO products( name, unit_price, units_in_stock, created_at , updated_at, category_id) VALUES ($1, $2, $3,$4,$5,$6) RETURNING id"
	var id int64
	err := p.db.QueryRow(query, product.Name, product.UnitPrice, product.UnitsInStock, product.CreatedAt, product.UpdatedAt, product.CategoryID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *ProductStorage) GetAllProduct() ([]Product, error) {
	query := "SELECT * FROM products"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.UnitPrice, &product.UnitsInStock, &product.CreatedAt, &product.UpdatedAt, &product.CategoryID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductStorage) GetOneProductByID(id int64) (Product, error) {
	query := "SELECT id,name, unit_price, units_in_stock, created_at, updated_at, category_id  FROM products WHERE id = $1"
	var product Product
	err := p.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.UnitPrice, &product.UnitsInStock, &product.CreatedAt, &product.UpdatedAt, &product.CategoryID)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (p *ProductStorage) UpdateProduct(product Product) (int64, error) {
	query := "UPDATE products SET name=$1,unit_price=$2,units_in_stock = $3,created_at = $4, updated_at =$5,category_id =$6 WHERE id = $7 RETURNING id"
	var id int64
	err := p.db.QueryRow(query, product.Name, product.UnitPrice, product.UnitsInStock, product.CreatedAt, product.UpdatedAt, product.CategoryID, product.ID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *ProductStorage) DeleteProductByID(id int64) error {
	query := "DELETE FROM products WHERE id=$1"
	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
