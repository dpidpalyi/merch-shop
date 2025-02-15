package dbinit

import (
	"context"
	"database/sql"
	"fmt"
	"merch-shop/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	fmt.Println(cfg.GetDSN())
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
