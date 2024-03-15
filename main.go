package main

import (
	"context"
	"database/sql"
	"donation-service/internal/services"
	"donation-service/internal/storage/listeners"
	"log"

	"github.com/vtopc/go-monobank"

	mono2 "donation-service/internal/services/mono"
	storage2 "donation-service/internal/storage"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/donation-service")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	logger := log.Default()

	donationsCache := storage2.NewCache()
	cacheListener := listeners.NewCacheListener(donationsCache)

	storage := storage2.NewStorage(db, cacheListener)

	token := "token"
	client := monobank.NewPersonalClient(nil).WithAuth(monobank.NewPersonalAuthorizer(token))
	mono := mono2.NewService(&client, storage, logger)

	service := services.NewDonationService(storage, logger, &mono, donationsCache)

	server := NewServer(service, logger)

	go mono.PoolAccountInfo(context.Background(), mono2.MinTs)

	server.Run()
}
