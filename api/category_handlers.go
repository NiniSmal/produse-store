package api

import (
	"Produse_store/storage"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CategoryHandler struct {
	storage *storage.CategoryStorage
}

func NewCategoryHandler(cs *storage.CategoryStorage) *CategoryHandler {
	return &CategoryHandler{
		storage: cs,
	}
}

func HandlerError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.Write([]byte("Sorry , error in the program"))
}

func (c *CategoryHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("home page")
}

func (c *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category storage.Category
	err := json.NewDecoder(r.Body).Decode(&category) //b, err := io.ReadAll(r.Body) ,err = json.Unmarshal(b, &task)
	if err != nil {
		HandlerError(w, err)
		return
	}
	err = category.Validate()
	if err != nil {
		HandlerError(w, err)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandlerError(w, err)
		return
	}
	category.ID, err = strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		HandlerError(w, err)
		return
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt
	_, err = c.storage.Save(category)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (c *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandlerError(w, err)
		return
	}
	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		HandlerError(w, err)
		return
	}
	categories, err := c.storage.GetAllCategories(userID)
	if err != nil {
		HandlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (c *CategoryHandler) GetOneCategoryByID(w http.ResponseWriter, r *http.Request) {
	idRequest := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idRequest)
	if err != nil {
		HandlerError(w, err)
		return
	}
	category, err := c.storage.GetCategoryByID(int64(id))
	if err != nil {
		HandlerError(w, err)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		HandlerError(w, err)
		return
	}
	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		HandlerError(w, err)
	}
	if category.UserID != userID {
		http.Error(w, "Not your category", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(category)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (c *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category storage.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		HandlerError(w, err)
		return
	}
	category.UpdatedAt = time.Now()
	_, err = c.storage.UpdateCategory(category)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (c *CategoryHandler) DeleteOneCategory(w http.ResponseWriter, r *http.Request) {
	idCat := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idCat)
	if err != nil {
		HandlerError(w, err)
	}
	err = c.storage.DeleteCategoryByID(int64(id))
	if err != nil {
		HandlerError(w, err)
	}
}
