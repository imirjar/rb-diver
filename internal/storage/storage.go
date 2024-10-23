package storage

import (
	"context"

	"github.com/imirjar/rb-diver/internal/storage/reports"
	"github.com/imirjar/rb-diver/internal/storage/target"
)

type Config interface {
	GetDiverTargetDB() string
}

type ReportsStore interface {
	GetQuery(ctx context.Context, id string) (string, error)
	GetAllReports(ctx context.Context) (string, error)
}

type Target interface {
	ExecuteQuery(ctx context.Context, query string) ([]map[string]any, error)
}

type Storage struct {
	ReportsStore
	Target
}

func New(conn string) *Storage {
	return &Storage{
		ReportsStore: reports.New(),
		Target:       target.New(conn),
	}
}
