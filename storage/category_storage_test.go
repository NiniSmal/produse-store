package storage

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestCategory_Validate(t *testing.T) {
	// поле заполнено
	category := Category{Name: "fish"}
	err := category.Validate()
	require.NoError(t, err)

	// длинное имя
	category = Category{Name: strings.Repeat("n", 100)}
	err = category.Validate()
	require.NoError(t, err)
}

func TestCategory_Validate_Error(t *testing.T) {
	// коротоке имя
	category := Category{Name: "dd"}
	err := category.Validate()
	require.Error(t, err)

	//длинное имя
	category = Category{Name: strings.Repeat("n", 101)}
	err = category.Validate()
	require.Error(t, err)
}

func TestCategoryStorage_Save(t *testing.T) {
	db := CategoryConnection(t)
	ct := NewCategoryStorage(db)

	category := Category{
		Name:   uuid.New().String(),
		UserID: 1,
	}
	id, err := ct.Save(category)
	require.NoError(t, err)

	dbCategory, err := ct.GetCategoryByID(id)
	require.NoError(t, err)
	require.NotEmpty(t, dbCategory)
}

func TestCategoryStorage_GetAllCategories(t *testing.T) {
	db := CategoryConnection(t)
	ct := NewCategoryStorage(db)
	userID := int64(1)
	createdAt := time.Now().Round(time.Millisecond)
	updatedAt := createdAt
	categories := []Category{
		{Name: uuid.New().String(), CreatedAt: createdAt, UpdatedAt: updatedAt, UserID: userID},
		{Name: uuid.New().String(), CreatedAt: createdAt, UpdatedAt: updatedAt, UserID: userID},
	}

	for index, category := range categories {
		id, err := ct.Save(category)
		require.NoError(t, err)
		categories[index].ID = id
	}
	dbCategory, err := ct.GetAllCategories(userID)
	require.NoError(t, err)
	require.Contains(t, dbCategory, categories)
}
func TestCategoryStorage_UpdateCategory(t *testing.T) {
	db := CategoryConnection(t)
	ct := NewCategoryStorage(db)

	category := Category{
		Name:   uuid.New().String(),
		UserID: 1,
	}
	id, err := ct.Save(category)
	require.NoError(t, err)

	createdAt := time.Now().Round(time.Millisecond)
	updatedAt := createdAt

	category1 := Category{
		ID:        id,
		Name:      uuid.NewString(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    1,
	}

	id, err = ct.UpdateCategory(category1)
	require.NoError(t, err)

	dbCategory, err := ct.GetCategoryByID(id)
	require.NoError(t, err)
	require.Equal(t, category1.ID, dbCategory.ID)
	require.Equal(t, category1.Name, dbCategory.Name)
	require.Equal(t, category1.UserID, dbCategory.UserID)
	require.Equal(t, category1.CreatedAt.Unix(), dbCategory.CreatedAt.Unix())
	require.Equal(t, category1.UpdatedAt.Unix(), dbCategory.UpdatedAt.Unix())
}

func TestCategoryStorage_GetCategoryByID_Error(t *testing.T) {
	db := CategoryConnection(t)
	ct := NewCategoryStorage(db)
	_, err := ct.GetCategoryByID(123)
	require.Error(t, err)
}
func TestCategoryStorage_DeleteCategoryByID(t *testing.T) {
	db := CategoryConnection(t)
	ct := NewCategoryStorage(db)
	category := Category{
		Name: uuid.NewString(), UserID: 1,
	}
	id, err := ct.Save(category)
	require.NoError(t, err)

	err = ct.DeleteCategoryByID(id)
	require.NoError(t, err)

	_, err = ct.GetCategoryByID(id)
	require.Error(t, err)

}

func CategoryConnection(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:8001/postgres?sslmode=disable")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = db.Close()
		require.NoError(t, err)
	})
	err = db.Ping()
	require.NoError(t, err)
	return db
}
