package psql

import (
	"github.com/jmoiron/sqlx"
	"github.com/walnuts1018/2024-AC-hacking/config"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(config config.Config) (*Client, error) {
	db, err := sqlx.Open("postgres", config.PSQLDSN)
	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) DB() *sqlx.DB {
	return c.db
}
