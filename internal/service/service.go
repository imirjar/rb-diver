package service

import (
	"context"

	"github.com/imirjar/rb-diver/internal/models"
)

type Service struct {
	Storage Storage
	RS      ReportsStorage
}

type ReportsStorage interface {
	GetReports(context.Context) ([]models.Report, error)
	GetReport(context.Context, string) (models.Report, error)
}

type Storage interface {
	ExecuteReport(context.Context, string) (models.Data, error)
}

func New() *Service {
	return &Service{}
}

func (s Service) ReportExecute(ctx context.Context, id string) (models.Data, error) {
	report, err := s.RS.GetReport(ctx, id)
	if err != nil {
		return models.Data{}, err
	}
	return s.Storage.ExecuteReport(ctx, report.Query)
}

func (s Service) ReportsList(ctx context.Context, roles []string) ([]models.Report, error) {
	return s.RS.GetReports(ctx)
}
