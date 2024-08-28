package settings

import (
	"fmt"
	"log"

	"github.com/ThailanTec/go-transactional-outbox/src/config"
	"github.com/jmoiron/sqlx"
)

func PostgresClient(cfg config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Println("Erro ao conectar no banco de dados:", err)
		return nil, err
	}

	return db, nil
}
