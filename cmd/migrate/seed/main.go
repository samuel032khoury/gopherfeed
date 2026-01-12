package main

import (
	"log"

	"github.com/samuel032khoury/gopherfeed/internal/db"
	"github.com/samuel032khoury/gopherfeed/internal/env"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

func main() {
	url := env.GetString("DB_URL", "postgres://user:password@localhost:5432/gopherfeed?sslmode=disable")
	conn, err := db.New(url, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewPostgresStorage(conn)
	seed(store)
}
