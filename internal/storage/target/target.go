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
	// if ping(dbConn) {
	pool, err := pgxpool.New(ctx, dbConn)
	if err != nil {
		panic(err)
	}
	return &TargetDB{
		pool: pool,
	}
	// } else {
	// 	log.Print(dbConn)
	// 	panic("YOU CAN'T RUN DIVER WITHOUT DB CONNECTION!!!")
	// }

}

// Execute selected report tranform report data to models.Data
func (t *TargetDB) ExecuteQuery(ctx context.Context, query string) (models.Data, error) {

	// DB request
	rows, err := t.pool.Query(ctx, query)
	if err != nil {
		log.Print(err)
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
			return models.Data{}, err
		}
		valueRows = append(valueRows, values)
	}

	// create answer
	results := models.Data{
		Columns: columns,
		Values:  valueRows,
	}

	// Получаем имена колонок
	// fieldDescriptions := rows.FieldDescriptions()
	// columns := make([]string, len(fieldDescriptions))
	// for i, fd := range fieldDescriptions {
	// 	columns[i] = string(fd.Name)
	// }

	// var results []map[string]interface{}
	// for rows.Next() {
	// 	values, err := rows.Values()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Создаем динамическую запись
	// 	record := make(map[string]interface{})
	// 	for i, col := range columns {
	// 		record[col] = values[i]
	// 	}
	// 	results = append(results, record)
	// }
	// log.Print(results)
	return results, nil
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
