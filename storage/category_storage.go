package storage

import (
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	UserID    int64     `json:"user_id"`
}

type CategoryUpdate struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryStorage struct {
	db *sql.DB
}

func NewCategoryStorage(c *sql.DB) *CategoryStorage {
	return &CategoryStorage{
		db: c,
	}
}
func (c *Category) Validate() error {
	r := utf8.RuneCountInString(c.Name)
	if r < 3 || r > 100 {
		return errors.New("the name must be more 3 and less 100")
	}
	return nil
}

func (c *CategoryStorage) Save(category Category) (int64, error) {
	query := "INSERT INTO categories(name, created_at, updated_at, user_id) VALUES ($1,$2,$3,$4) RETURNING id"
	var id int64
	err := c.db.QueryRow(query, category.Name, category.CreatedAt, category.UpdatedAt, category.UserID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) GetAllCategories(userID int64) ([]Category, error) {
	query := "SELECT * FROM categories WHERE user_id =$1"
	rows, err := c.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.UserID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryStorage) GetCategoryByID(id int64) (Category, error) {
	query := "SELECT id, name, created_at, updated_at, user_id FROM categories WHERE id = $1"
	var category Category
	err := c.db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.UserID)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func (c *CategoryStorage) UpdateCategory(category Category) (int64, error) {
	query := "UPDATE categories SET name = $1,created_at=$2, updated_at = $3,user_id = $4 WHERE id = $5 RETURNING id"
	var id int64
	err := c.db.QueryRow(query, category.Name, category.CreatedAt, category.UpdatedAt, category.UserID, category.ID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CategoryStorage) DeleteCategoryByID(id int64) error {
	query := "DELETE FROM categories  WHERE id = $1"
	_, err := c.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
