package db

import (
	"fmt"
	"log"

	"github.com/tokopedia/panics"
	"github.com/tokopedia/sqlt"
)

var configDB = struct {
	host     string
	port     string
	user     string
	password string
}{
	host:     "devel-postgre.tkpd",
	port:     "5432",
	user:     "yp180102",
	password: "OiQr2df1DSy4Q0",
}

var DB dbconnection

type dbconnection struct {
	DBTDev   *sqlt.DB
	DBTOrder *sqlt.DB
}

func InitDB() {
	//shop DB
	dbConnTDev := generateDBConn("tokopedia-dev-db")

	// product DB
	dbConnTOrder := generateDBConn("tokopedia-order")

	//store all DB Connection
	DB = dbconnection{
		DBTDev:   dbConnTDev,
		DBTOrder: dbConnTOrder,
	}
}

func newPostgresDB(dsn string) *sqlt.DB {

	db, err := sqlt.Open("postgres", dsn)
	if err != nil {
		panics.Capture("Cant connect to DB ", dsn)
		log.Fatalf("failed to connect to %s: %s", dsn, err.Error())
	}

	if err := db.Ping(); err != nil {
		panics.Capture("Can't reach DB ", dsn)
		log.Fatalf("%s database is unreachable: %s", dsn, err.Error())
	}

	db.SetMaxOpenConnections(200)

	return db
}

func generateDBConn(dbName string) *sqlt.DB {
	temp_con := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", configDB.host, configDB.port, configDB.user, configDB.password, dbName)
	dbConn := newPostgresDB(temp_con)
	return dbConn
}
