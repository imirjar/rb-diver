package service

import (
	"context"
	"log"

	"github.com/imirjar/rb-diver/internal/models"
)

type Service struct {
	Storage Storage
}

type Storage interface {
	GetReport(context.Context, string) (models.Report, error)
	GetReports(context.Context) ([]models.Report, error)
	GetReportsByRole(context.Context, []string) ([]models.Report, error)

	GetRoles(ctx context.Context) ([]models.Role, error)
	GetRolesByReportID(context.Context, string) ([]models.Role, error)

	ExecuteReport(context.Context, string) (models.Data, error)
}

func New() *Service {
	return &Service{}
}

func (s Service) ReportExecute(ctx context.Context, id string) (models.Data, error) {
	report, err := s.Storage.GetReport(ctx, id)
	if err != nil {
		log.Print("GET QUERY ERROR:", err)
		return models.Data{}, err
	}

	data, err := s.Storage.ExecuteReport(ctx, report.Query)
	if err != nil {
		log.Print("EXECUTE QUERY ERROR:", err)
		return models.Data{}, err
	}

	// log.Println(query, string(data.Raw))
	return data, nil
}

func (s Service) ReportInfo(ctx context.Context, id string) (models.Report, error) {
	query, err := s.Storage.GetReport(ctx, id)
	if err != nil {
		log.Print("GET QUERY ERROR:", err)
		return models.Report{}, err
	}

	// log.Println(query, string(data.Raw))
	return query, nil
}

func (s Service) ReportsList(ctx context.Context, roles []string) ([]models.Report, error) {
	if len(roles) > 0 {
		log.Print(roles)
		// rolesStrIDs := []string{}
		// rolesIDS, err := s.Storage.GetRoles()
		return s.Storage.GetReportsByRole(ctx, roles)
	}
	return s.Storage.GetReports(ctx)
}

func (s Service) RoleList(ctx context.Context, repID string) ([]models.Role, error) {
	if repID != "" {
		return s.Storage.GetRolesByReportID(ctx, repID)
	} else {
		// log.Print("no report")
		return s.Storage.GetRoles(ctx)
	}
}
