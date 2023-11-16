package main

import (
	"Produse_store/api"
	"Produse_store/storage"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	//docker run  -d -p 8001:5432 -e POSTGRES_PASSWORD=dev -e POSTGRES_DATABASE=postgres postgres

	db, err := sql.Open("postgres", "postgres://postgres:dev@localhost:8001/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db - OK")

	cs := storage.NewCategoryStorage(db)
	ch := api.NewCategoryHandler(cs)
	ps := storage.NewProductStorage(db)
	ph := api.NewProductHandler(ps)
	us := storage.NewUserStorage(db)
	uh := api.NewUserHandler(us)

	router := http.NewServeMux()
	router.HandleFunc("/", ch.HomePage)
	router.HandleFunc("/categories/create", ch.CreateCategory)
	router.HandleFunc("/categories/getList", ch.GetAllCategories)
	router.HandleFunc("/categories/getOn", ch.GetOneCategoryByID)
	router.HandleFunc("/categories/update", ch.UpdateCategory)
	router.HandleFunc("/categories/deleteOne", ch.DeleteOneCategory)

	router.HandleFunc("/products/create", ph.CreateProduct)
	router.HandleFunc("/products/getList", ph.GetAllProduct)
	router.HandleFunc("/products/getOn", ph.GetOneProductByID)
	router.HandleFunc("/products/update", ph.UpdateProduct)
	router.HandleFunc("/products/delete", ph.DeleteProduct)

	router.HandleFunc("/users/create", uh.CreateUser)
	router.HandleFunc("/users/getOn", uh.GetUserByID)
	router.HandleFunc("/auth/login", uh.Login)

	err = http.ListenAndServe(":8093", router)
	if err != nil {
		log.Fatal(err)
	}
}
