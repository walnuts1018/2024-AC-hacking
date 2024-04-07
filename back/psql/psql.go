package psql

import (
	"fmt"

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

func (c *Client) Init() error {
	if err := c.CreateUserTable(); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateUserTable() error {
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT, password TEXT)")
	return err
}

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func (c *Client) Login(username, password string) (User, error) {
	var user User
	err := c.db.Get(&user, fmt.Sprintf("SELECT * FROM users WHERE username = '%s' AND password = '%s'", username, password))
	return user, err
}

func (c *Client) Register(username, password string) error {

	if _, err := c.GetUser(username); err == nil {
		return fmt.Errorf("user already exists")
	}

	_, err := c.db.Exec(fmt.Sprintf("INSERT INTO users (username, password) VALUES ('%s', '%s')", username, password))
	return err
}

func (c *Client) GetUser(username string) ([]User, error) {
	var users []User
	err := c.db.Get(&users, fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username))
	return users, err
}

func (c *Client) ChangePassword(username, oldPassword, password string) error {
	user, err := c.Login(username, oldPassword)
	if err != nil {
		return err
	}

	_, err = c.db.Exec(fmt.Sprintf("UPDATE users SET password = '%s' WHERE username = '%s'", password, user.Username))
	return err
}
