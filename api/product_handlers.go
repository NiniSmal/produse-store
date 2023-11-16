package api

import (
	"Produse_store/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ProductHandler struct {
	storage *storage.ProductStorage
}

func NewProductHandler(ps *storage.ProductStorage) *ProductHandler {
	return &ProductHandler{
		storage: ps,
	}
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product storage.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		HandlerError(w, err)
		return
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt
	_, err = p.storage.Save(product)
	if err != nil {
		HandlerError(w, err)
		return
	}
}

func (p *ProductHandler) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := p.storage.GetAllProduct()
	if err != nil {
		HandlerError(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (p *ProductHandler) GetOneProductByID(w http.ResponseWriter, r *http.Request) {
	idPr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idPr)
	if err != nil {
		HandlerError(w, err)
		return
	}
	product, err := p.storage.GetOneProductByID(int64(id))
	if err != nil {
		HandlerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product storage.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		HandlerError(w, err)
		return
	}
	product.UpdatedAt = time.Now()

	_, err = p.storage.UpdateProduct(product)
	if err != nil {
		HandlerError(w, err)
		return
	}
}
func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idPr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idPr)
	if err != nil {
		HandlerError(w, err)
	}
	err = p.storage.DeleteProductByID(int64(id))
	if err != nil {
		HandlerError(w, err)
		return
	}
}
