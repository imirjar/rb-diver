package reports

import (
	"context"
	"encoding/json"
	"log"

	"github.com/imirjar/rb-diver/internal/models"
	_ "modernc.org/sqlite"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportsStore struct {
	mongoClient *mongo.Client
}

func New() *ReportsStore {

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	store := ReportsStore{
		mongoClient: mongoClient,
	}

	// if err := store.Ping(); err != nil {
	// 	panic(err)
	// }

	return &store
}

func (r *ReportsStore) Ping() error {
	return nil
}

// func (r *ReportsStore) GetQuery(ctx context.Context, id string) (string, error) {

// 	var data string

// 	row := r.dbConn.QueryRowContext(ctx, "SELECT query FROM reports WHERE id=$1;", id)
// 	err := row.Scan(&data)
// 	if err != nil {
// 		return err.Error(), err
// 	}

// 	return data, nil
// }

func (r *ReportsStore) GetQuery(ctx context.Context, id string) (string, error) {

	collection := r.mongoClient.Database("reports").Collection("reports")
	filter := bson.D{{"_id", id}}

	report := models.Report{}

	if err := collection.FindOne(ctx, filter).Decode(&report); err != nil {
		return err.Error(), err
	}

	return report.Query, nil
}

func (r *ReportsStore) GetAllReports(ctx context.Context) (string, error) {

	rows, err := r.dbConn.QueryContext(ctx, "SELECT * FROM reports;")
	if err != nil {
		return err.Error(), err
	}

	var reports []models.Report
	for rows.Next() {

		var rep models.Report
		err = rows.Scan(&rep.Id, &rep.Name, &rep.Query)
		if err != nil {
			return err.Error(), err
		}
		reports = append(reports, rep)
	}
	result, err := json.Marshal(&reports)
	if err != nil {
		return err.Error(), err
	}

	return string(result), nil
}
