package reports

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/imirjar/rb-diver/internal/models"
	_ "modernc.org/sqlite"
)

type ReportsStore struct {
	dbConn *sql.Conn
}

func New() *ReportsStore {
	db, err := sql.Open("sqlite", "db/reports")
	if err != nil {
		// if err.Error() == "" {
		// 	log.Print(err.Error())
		// }
		log.Println("-1->", err.Error(), "<--")
		// panic(err)
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Println("-2->", err.Error(), "<--")
		panic(err)
	}

	store := ReportsStore{
		dbConn: conn,
	}

	// if err := store.Ping(); err != nil {
	// 	panic(err)
	// }

	return &store
}

func (r *ReportsStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return r.dbConn.PingContext(ctx)
}

func (r *ReportsStore) GetQuery(ctx context.Context, id string) (string, error) {

	var data string

	row := r.dbConn.QueryRowContext(ctx, "SELECT query FROM reports WHERE id=$1;", id)
	err := row.Scan(&data)
	if err != nil {
		return err.Error(), err
	}

	return data, nil
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
