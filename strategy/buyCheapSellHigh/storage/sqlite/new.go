package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	setupCode string
	db        *sql.DB
}

var db *sql.DB

func New(path string, setupCode string) (*Storage, error) {
	if db != nil {
		return initStorage(setupCode)
	}

	err := error(nil)
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return initStorage(setupCode)
}

func initStorage(setupCode string) (*Storage, error) {

	storage := Storage{
		db:        db,
		setupCode: setupCode,
	}

	orderTableName := getOrderTableName(setupCode)

	initQueries := []string{
		`CREATE TABLE IF NOT EXISTS ` + orderTableName + ` (
			Id                        INTEGER PRIMARY KEY AUTOINCREMENT,
			DomainOrderId             TEXT,
			FilledPrice               REAL,
			FilledQty                 REAL,
			Side                      TEXT,
			CreatedTime               INTEGER,
			UpdatedTime               INTEGER,
			MainCurrencyAmountBefore  REAL,
			TradeCurrencyAmountBefore REAL,
			AveragePrice              REAL, 
			UNIQUE (DomainOrderId)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_CreatedTime ON ` + orderTableName + ` (CreatedTime)`,
		`CREATE INDEX IF NOT EXISTS idx_UpdatedTime ON ` + orderTableName + ` (UpdatedTime)`,
	}

	for _, query := range initQueries {
		_, err := storage.db.Exec(query)

		if err != nil {
			return nil, err
		}
	}

	return &storage, nil
}

func getOrderTableName(setupCode string) string {
	return `orders_` + setupCode
}
