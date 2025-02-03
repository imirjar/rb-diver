package storage

import (
	"context"

	"github.com/imirjar/rb-diver/internal/models"
	"github.com/imirjar/rb-diver/internal/storage/self"
	"github.com/imirjar/rb-diver/internal/storage/target"
)

type Config interface {
	GetDiverTargetDB() string
}

type SelfStore interface {
	GetReport(context.Context, string) (models.Report, error)
	GetReports(context.Context) ([]models.Report, error)
	GetReportsByRole(context.Context, []string) ([]models.Report, error)

	GetRoles(context.Context) ([]models.Role, error)
	GetRolesByReportID(context.Context, string) ([]models.Role, error)
}

type TargetStore interface {
	ExecuteReport(ctx context.Context, query string) (models.Data, error)
}

type Storage struct {
	SelfStore
	TargetStore
}

func New(conn string) *Storage {
	return &Storage{
		SelfStore:   self.New(),
		TargetStore: target.New(context.Background(), conn),
	}
}
