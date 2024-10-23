package target

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config interface {
	GetDiverTargetDB() string
}

type TargetDB struct {
	pool *sql.DB
}

func New(dbconn string) *TargetDB {
	pool, err := sql.Open("pgx", dbconn)
	if err != nil {
		panic(err)
	}
	return &TargetDB{
		pool: pool,
	}
}

func (t *TargetDB) ExecuteQuery(ctx context.Context, query string) ([]map[string]any, error) {
	rows, err := t.pool.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var allMaps []map[string]any
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		resultMap := make(map[string]interface{})
		for i, val := range values {
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}

	return allMaps, nil
}
