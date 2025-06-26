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

func New() *TargetDB {
	return &TargetDB{}
}

func (t *TargetDB) Connect(ctx context.Context, dbConn string) error {
	pool, err := pgxpool.New(ctx, dbConn)
	if err != nil {
		return err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Println("Внимание! Подключение к базе данных отсутствует!")
		return err
	}
	t.pool = pool

	return nil
}

func (t *TargetDB) Disconnect() {
	t.pool.Close()
}

// Execute selected report tranform report data to models.Data
func (t *TargetDB) ExecuteReport(ctx context.Context, query string) (models.Data, error) {

	log.Print(query)
	// DB request
	rows, err := t.pool.Query(ctx, query)
	if err != nil {
		log.Printf("query %s with err %s", query, err)
		return models.Data{}, err
	}

	// future models.Data columns
	var columns = []string{}
	var valueRows = [][]any{}

	// DB table columns to models.Data columns
	fieldDescriptions := rows.FieldDescriptions()
	for _, v := range fieldDescriptions {
		columns = append(columns, v.Name)
	}

	// DB table values to models.Data columns list
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Print("EOF!!!!!")
			return models.Data{}, err
		}
		valueRows = append(valueRows, values)
	}

	// create answer
	results := models.Data{
		Columns: columns,
		Values:  valueRows,
	}

	return results, nil
}
