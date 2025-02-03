package self

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type SelfStore struct {
	dbConn *sql.Conn
}

func New() *SelfStore {

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

	store := SelfStore{
		dbConn: conn,
	}

	// if err := store.Ping(); err != nil {
	// 	panic(err)
	// }

	return &store
}

func (r *SelfStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return r.dbConn.PingContext(ctx)
}
