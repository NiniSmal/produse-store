package api

import (
	"Produse_store/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	storage *storage.UserStorage
}

func NewUserHandler(u *storage.UserStorage) *UserHandler {
	return &UserHandler{
		storage: u,
	}
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user storage.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}

	_, err = u.storage.Save(user)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idUser := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idUser)
	if err != nil {
		HandlerError(w, err)
		return
	}

	user, err := u.storage.GetUserByID(int64(id))
	if err != nil {
		HandlerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "json/application")

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user storage.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HandlerError(w, err)
		return
	}
	id, err := u.storage.Login(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   strconv.FormatInt(id, 10),
		Path:    "/",
		MaxAge:  3600,
		Expires: time.Now().Add(time.Hour),
	}
	http.SetCookie(w, cookie)
}
