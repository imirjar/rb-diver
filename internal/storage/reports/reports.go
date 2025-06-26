package reports

import (
	"context"
	"fmt"
	"log"

	"github.com/imirjar/rb-diver/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Storage struct {
	client  *mongo.Client
	reports []models.Report
}

func New() *Storage {
	// store := &Storage{
	// 	reports: []models.Report{
	// 		models.Report{
	// 			Id:    "1",
	// 			Name:  "first",
	// 			Query: "SELECT * FROM db;",
	// 		},
	// 		models.Report{
	// 			Id:    "2",
	// 			Name:  "second",
	// 			Query: "SELECT * FROM db;",
	// 		},
	// 		models.Report{
	// 			Id:    "3",
	// 			Name:  "third",
	// 			Query: "SELECT * FROM db;",
	// 		},
	// 	},
	// }
	store := &Storage{}
	return store
}

func (s *Storage) Connect(ctx context.Context, conn string) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	// client, err := mongo.Connect(conn)
	if err != nil {
		panic(err)
	}
	s.client = client

	return nil
}

func (s *Storage) Disconnect() error {
	log.Print("mongo disconn ok")
	return s.client.Disconnect(context.Background())
}

func (s *Storage) GetReports(ctx context.Context) ([]models.Report, error) {
	collection := s.client.Database("PoliglotimCourses").Collection("reports")

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Print("ERROR DiverReports ", err)
		return nil, fmt.Errorf("failed to find reports: %v", err)
	}
	defer cur.Close(ctx)

	var reports []models.Report

	for cur.Next(ctx) {
		var result models.Report
		if err := cur.Decode(&result); err != nil {
			log.Printf("ERROR decoding report: %v", err)
			continue // или return nil, err если хотите прервать при ошибке
		}
		reports = append(reports, result)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return reports, nil
}

func (s *Storage) GetReport(ctx context.Context, repID string) (models.Report, error) {
	collection := s.client.Database("PoliglotimCourses").Collection("reports")

	// Преобразуем строковый ID в ObjectID
	objID, err := primitive.ObjectIDFromHex(repID)
	if err != nil {
		return models.Report{}, fmt.Errorf("invalid report ID format: %v", err)
	}

	var report models.Report
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Report{}, fmt.Errorf("report not found")
		}
		log.Printf("ERROR getting report: %v", err)
		return models.Report{}, fmt.Errorf("failed to get report: %v", err)
	}

	return report, nil
}
