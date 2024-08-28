package repository

import (
	"fmt"
	"time"

	"github.com/ThailanTec/go-transactional-outbox/domain"
	"github.com/jmoiron/sqlx"
)

type OrderService struct {
	db *sqlx.DB
}

func NewServiceOrder(db *sqlx.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) CreateOrder(order domain.Order) (int, error) {
	query := "INSERT INTO orders (item, quantity, created_at) VALUES ($1, $2, $3) RETURNING id"

	var orderID int
	err := s.db.QueryRow(query, order.Item, order.Quantity, time.Now()).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (s *OrderService) CreateOutboxMessage(orderOutbox domain.OutboxMessage) error {
	query := "INSERT INTO outbox (payload, processed, created_at) VALUES ($1, $2, $3)"

	_, err := s.db.Exec(query, orderOutbox.Payload, false, time.Now())
	if err != nil {

		return err
	}

	return nil
}

func (s *OrderService) GetOutboxMessage(id int) (out *domain.OutboxMessage, e error) {
	query := "SELECT id, payload, processed FROM outbox WHERE id = $1 AND processed = false"
	row := s.db.QueryRow(query, id)
	var outbox domain.OutboxMessage

	err := row.Scan(&outbox.ID, &outbox.Payload, &outbox.Processed)
	if err != nil {
		fmt.Println("Errou aqui")
		return nil, err
	}

	return &outbox, nil
}

func (s *OrderService) UpdateOutboxMessage(id int) error {
	query := "UPDATE outbox SET processed = true WHERE id = $1"

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
