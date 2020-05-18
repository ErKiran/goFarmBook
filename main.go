package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ErKiran/node/db"
	"github.com/ErKiran/node/domain"
	"github.com/ErKiran/node/handlers"
	"github.com/go-pg/pg/v10"
)

func main() {
	DB := db.New(&pg.Options{
		User:     "golang",
		Password: "func",
		Database: "node",
	})

	defer DB.Close()

	domainDB := domain.DB{
		UserRepo: db.NewUserRepo(DB),
	}

	d := &domain.Domain{DB: domainDB}

	r := handlers.SetUpRouter(d)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("cannot start server %v", err)
	}
	fmt.Println("Server is starting.... Tada")
}
