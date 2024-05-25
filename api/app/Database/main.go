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

	if err := database.checkAndCreateTable(); err != nil {
		return nil, err
	}

	return &database, nil
}

func (ctr *API) checkAndCreateTable() error {
	var table string
	err := ctr.db.QueryRow("SHOW TABLES LIKE 'ThirdApp'").Scan(&table)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if table == "" {
		_, err := ctr.db.Exec(`
		CREATE TABLE ThirdApp (
			id              CHAR(36),
			name            VARCHAR(255)     NOT NULL,
			client          INT             NOT NULL,
			salt            VARCHAR(64)     NOT NULL          COMMENT 'secret salt',
			hash            VARCHAR(64)     NOT NULL          COMMENT 'secret hash',
			callback        VARCHAR(255)    NOT NULL,
			description     TEXT            DEFAULT ''        COMMENT 'app description',
			permissions     INT             NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (client) REFERENCES User(id)
		)`)
		if err != nil {
			return err
		}
	}

	return nil
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
