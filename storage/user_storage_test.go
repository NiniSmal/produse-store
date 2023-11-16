package storage

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserStorage_Save(t *testing.T) {
	db := UserConnection(t)
	ut := NewUserStorage(db)

	user := User{
		Login:    uuid.New().String(),
		Password: uuid.New().String(),
	}
	id, err := ut.Save(user)
	require.NoError(t, err)

	dbUser, err := ut.GetUserByID(id)
	require.NoError(t, err)
	require.NotEmpty(t, dbUser.ID)
	require.Equal(t, user.Login, dbUser.Login)
	require.Empty(t, dbUser.Password)
}
func TestUserStorage_GetUserByID_Error(t *testing.T) {
	db := UserConnection(t)

	us := NewUserStorage(db)
	_, err := us.GetUserByID(123)
	require.Error(t, err)
}
func TestUserStorage_Login(t *testing.T) {

	db := UserConnection(t)
	us := NewUserStorage(db)

	user := User{
		Login:    uuid.New().String(),
		Password: uuid.New().String(),
	}
	id, err := us.Save(user)
	require.NoError(t, err)

	dbId, err := us.Login(user.Login, user.Password)
	require.NoError(t, err)
	require.Equal(t, id, dbId)

}

func UserConnection(t *testing.T) *sql.DB {
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
