package storage

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestProductStorage_Save(t *testing.T) {
	db := ProductConnection(t)
	pt := NewProductStorage(db)

	product := Product{
		Name:         uuid.NewString(),
		UnitPrice:    2,
		UnitsInStock: 10,
		CategoryID:   1,
	}
	id, err := pt.Save(product)
	require.NoError(t, err)

	dbProduct, err := pt.GetOneProductByID(id)
	require.NoError(t, err)
	require.NotEmpty(t, dbProduct)
}

func TestProductStorage_GetOneProductByID(t *testing.T) {
	db := ProductConnection(t)
	pt := NewProductStorage(db)

	_, err := pt.GetOneProductByID(123)
	require.Error(t, err)
}
func TestProductStorage_UpdateProduct(t *testing.T) {
	db := ProductConnection(t)
	pt := NewProductStorage(db)

	createdAt := time.Now().Round(time.Millisecond)
	updatedAt := createdAt

	product := Product{
		Name:         uuid.NewString(),
		UnitPrice:    0,
		UnitsInStock: 0,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		CategoryID:   1,
	}
	id, err := pt.Save(product)
	require.NoError(t, err)

	product1 := Product{
		ID:           id,
		Name:         uuid.NewString(),
		UnitPrice:    0,
		UnitsInStock: 0,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		CategoryID:   1,
	}
	id, err = pt.UpdateProduct(product1)
	require.NoError(t, err)
	bdProduct, err := pt.GetOneProductByID(id)
	require.NoError(t, err)
	require.Equal(t, product1.ID, bdProduct.ID)
	require.Equal(t, product1.Name, bdProduct.Name)
	require.Equal(t, product1.CategoryID, bdProduct.CategoryID)
	require.Equal(t, product1.CreatedAt.Unix(), bdProduct.CreatedAt.Unix())
	require.Equal(t, product1.UpdatedAt.Unix(), bdProduct.UpdatedAt.Unix())

}

func TestProductStorage_DeleteProductByID(t *testing.T) {
	db := ProductConnection(t)
	pt := NewProductStorage(db)

	product := Product{
		Name:         "fish",
		UnitPrice:    0,
		UnitsInStock: 0,
		CategoryID:   1,
	}
	id, err := pt.Save(product)
	require.NoError(t, err)

	err = pt.DeleteProductByID(id)
	require.NoError(t, err)

	_, err = pt.GetOneProductByID(id)
	require.Error(t, err)
}
func ProductConnection(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:8001/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	t.Cleanup(func() {
		err = db.Close()
		require.NoError(t, err)
	})
	err = db.Ping()
	require.NoError(t, err)
	return db
}
