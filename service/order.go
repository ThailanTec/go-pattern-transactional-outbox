package service

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ThailanTec/go-transactional-outbox/domain"
	"github.com/ThailanTec/go-transactional-outbox/repository"
)

type OrderService interface {
	CreateOrder(order domain.Order) error
	TestRoutine(wg *sync.WaitGroup)
}

type orderService struct {
	orderRep *repository.OrderService
}

func NewOrderService(orderRep *repository.OrderService) *orderService {
	return &orderService{
		orderRep: orderRep,
	}
}

func (or *orderService) CreateOrder(order domain.Order) (int, error) {
	id, err := or.orderRep.CreateOrder(order)
	if err != nil {
		return 0, err
	}

	outbox := domain.OutboxMessage{
		Payload: order.Item,
	}

	err = or.orderRep.CreateOutboxMessage(outbox)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (or *orderService) TestRoutine(ctx context.Context, wg *sync.WaitGroup, id <-chan int, out *domain.OutboxMessage) {
	defer wg.Done()

	for {
		select {
		case ids, ok := <-id:
			if !ok {
				return
			}

			err := or.orderRep.UpdateOutboxMessage(ids)
			if err != nil {
				s := fmt.Sprintf("Erro ao atualizar outbox message para o ID %d", ids)
				fmt.Println(s)
			}

			s := repository.FicQueue(out)

			fmt.Println(s)
		case <-ctx.Done():
			log.Println("Contexto cancelado, encerrando goroutine")
			return
		}
	}
}
