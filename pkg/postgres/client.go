package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// newDB initialises a new db connection
func newDB(host, port, username, password, dbname string) (*sql.DB, error) {

	dnsStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s connect_timeout=3", host, port, username, password, dbname)

	db, err := sql.Open("postgres", dnsStr)
	if err != nil {
		fmt.Printf("sql.Open Err: %v\n", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping Err:", err)
		return nil, err
	}
	return db, nil
}

// Client is the client for the weatherapi database
type Client struct {
	database *sql.DB
}

// NewClient creates a new database client
func NewClient(host, port, username, password, dbName string) (*Client, error) {
	db, err := newDB(host, port, username, password, dbName)
	if err != nil {
		return nil, err
	}

	return &Client{db}, nil
}
