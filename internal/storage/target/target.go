package target

import (
	"context"
	"log"

	"github.com/imirjar/rb-diver/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config interface {
	GetDiverTargetDB() string
}

type TargetDB struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dbConn string) *TargetDB {
	if ping(dbConn) {
		pool, err := pgxpool.New(ctx, dbConn)
		if err != nil {
			panic(err)
		}
		return &TargetDB{
			pool: pool,
		}
	} else {
		panic("AHTUNG!!!")
	}

}

func (t *TargetDB) ExecuteQuery(ctx context.Context, query string) (*models.Data, error) {
	var data models.Data
	err := t.pool.QueryRow(ctx, query).Scan(&data.Raw)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &data, nil
}

func ping(conn string) bool {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conn)
	if err != nil {
		log.Print(err.Error())
		return false
	}
	if err = pool.Ping(ctx); err != nil {
		log.Println("Внимание! Подключение к базе данных отсутствует!")
		return false
	} else {
		return true
	}
}
