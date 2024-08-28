package repository

import (
	"fmt"

	"github.com/ThailanTec/go-transactional-outbox/domain"
)

func FicQueue(out *domain.OutboxMessage) string {
	str := fmt.Sprintf("Recebi os dados na Fila! ID: %d, payload: %s, processed: %t", out.ID, out.Payload, out.Processed)

	return str
}
