package orm

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	sw "single-window/pkg/postgres"
)

func NewPostgresORM(conn *sw.Postgres) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn: conn.DB,
		}),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init gorm, err: %w", err)
	}

	return db, nil
}
