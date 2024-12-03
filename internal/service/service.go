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
	GetQuery(context.Context, string) (string, error)
	ExecuteQuery(context.Context, string) (models.Data, error)
	ExecuteQueryMap(context.Context, string) ([]map[string]interface{}, error)
	GetAllReports(ctx context.Context) (string, error)
}

func New() *Service {
	return &Service{}
}

func (s Service) Execute(ctx context.Context, id string) (models.Data, error) {
	query, err := s.Storage.GetQuery(ctx, id)
	if err != nil {
		log.Print("GET QUERY ERROR:", err)
		return models.Data{}, err
	}

	data, err := s.Storage.ExecuteQuery(ctx, query)
	if err != nil {
		log.Print("EXECUTE QUERY ERROR:", err)
		return models.Data{}, err
	}

	// log.Println(query, string(data.Raw))
	return data, nil
}

func (s Service) ExecuteMap(ctx context.Context, id string) ([]map[string]interface{}, error) {
	query, err := s.Storage.GetQuery(ctx, id)
	if err != nil {
		log.Print("GET QUERY ERROR:", err)
		return nil, err
	}

	data, err := s.Storage.ExecuteQueryMap(ctx, query)
	if err != nil {
		log.Print("EXECUTE QUERY ERROR:", err)
		return nil, err
	}

	// log.Println(query, string(data.Raw))
	return data, nil
}

func (s Service) ReportsList(ctx context.Context) (string, error) {
	return s.Storage.GetAllReports(ctx)
}
