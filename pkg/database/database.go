package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
)

type Database struct {
	DB                    *sqlx.DB
	url                   string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	LogQueries            bool
}

type Options struct {
	URL                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	LogQueries            bool
}

func init() {
	driver := &sqlite3.SQLiteDriver{}
	sql.Register("sqlite3_extended", driver)
	sql.Register("sqlite3_extended_with_logs", sqlhooks.Wrap(driver, QueryLogger()))
}

// New with the given options.
// If no logger is provided, logs are discarded.
func New(opts Options) *Database {

	// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	opts.URL += "?_journal=WAL&_timeout=5000&_fk=true"

	return &Database{
		url:                   opts.URL,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		LogQueries:            opts.LogQueries,
	}
}

func (d *Database) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var driverName string
	if d.LogQueries {
		driverName = "sqlite3_extended_with_logs"
	} else {
		driverName = "sqlite3_extended"
	}

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, driverName, d.url)
	if err != nil {
		return err
	}

	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	return nil
}

type Transaction struct {
	Tx *sqlx.Tx
}

func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Transaction, error) {
	tx, err := d.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Transaction{Tx: tx}, nil
}

func (t *Transaction) Commit() error {
	return t.Tx.Commit()
}

func (t *Transaction) Rollback() error {
	return t.Tx.Rollback()
}
