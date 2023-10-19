package db

import (
	"os"

	"github.com/joho/godotenv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"oursos.com/packages/util"
)

func Connection() (*sql.DB, error) {
	err := godotenv.Load()

	util.CheckError(err)

	connStr := os.Getenv("CONN_DB")
	db, err := sql.Open("postgres", connStr)
	util.CheckError(err)
	return db, err
}
