package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ThailanTec/go-transactional-outbox/domain"
	"github.com/ThailanTec/go-transactional-outbox/repository"
	"github.com/ThailanTec/go-transactional-outbox/service"
	"github.com/ThailanTec/go-transactional-outbox/settings"
	"github.com/ThailanTec/go-transactional-outbox/src/config"
)

func main() {
	cfg := config.LoadConfig()

	db, err := settings.PostgresClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = settings.Migrations(db)
	if err != nil {
		log.Fatalf("Failed to generate a migrations: %v", err)
	}

	repo := repository.NewServiceOrder(db)
	service := service.NewOrderService(repo)

	Odata := domain.Order{
		Item:     "Redmi Note 12",
		Quantity: 1,
	}

	id, err := service.CreateOrder(Odata)
	if err != nil {
		log.Fatal(err)
	}

	out, err := repo.GetOutboxMessage(id)
	if err != nil {
		log.Fatal(err)
	}

	// Trabalhando com as Go routines
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go service.TestRoutine(context.Background(), &wg, ch, out)
	}

	go func() {
		select {
		case ch <- id:
			log.Print("Recebido com sucesso ", id)
		default:
			log.Print("Erro ao receber o ID ", id)
		}

		time.Sleep(10 * time.Second)
	}()

	wg.Wait()
}
