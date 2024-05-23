package Database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type API struct {
	db *sql.DB
}

type Config struct {
	SourceName string
}

func New(config Config) (*API, error) {
	database := API{}
	connect, err := sql.Open("mysql", config.SourceName)
	if err != nil {
		return nil, err
	}
	connect.SetMaxOpenConns(100)
	connect.SetMaxIdleConns(10)
	database.db = connect

	return &database, nil
}

func (api *API) UUID() (string, error) {
	connect := api.Connect()
	rows, err := connect.Query("SELECT UUID()")
	if err != nil {
		return "", err
	}
	rows.Next()

	var uuid string
	if err := rows.Scan(&uuid); err != nil {
		return "", err
	}

	return uuid, nil
}

func (api *API) Connect() *sql.DB {
	return api.db
}

func (api *API) Close() error {
	return api.db.Close()
}
