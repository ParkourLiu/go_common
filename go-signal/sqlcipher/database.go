package sqlcipher

import (
	"database/sql"
	"fmt"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

type DB struct {
	db *sql.DB
}

func Open(dbPath string, keyHex string) (*DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", dbPath, keyHex))
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func (c *DB) Create() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		name TEXT PRIMARY KEY NOT NULL,
		keys TEXT UNIQUE NOT NULL
	);`
	if _, err := c.db.Exec(createTableSQL); err != nil {
		return err
	}
	return nil
}
func (c *DB) Insert(name string, keys string) error {
	// 5. 插入数据
	insertSQL := "INSERT INTO users (name, keys) VALUES (?, ?);"
	_, err := c.db.Exec(insertSQL, name, keys)
	if err != nil {
		return err
	}
	return nil
}

func (c *DB) Update(name string, keys string) error {
	// 5. 插入数据
	insertSQL := "UPDATE users SET keys = ? WHERE name = ?;"
	_, err := c.db.Exec(insertSQL, keys, name)
	if err != nil {
		return err
	}
	return nil
}

func (c *DB) Query(name string) (string, string, error) {
	var keys string
	querySQL := "SELECT keys FROM users WHERE name = ?;"
	err := c.db.QueryRow(querySQL, name).Scan(&keys)
	if err != nil {
		return "", "", err
	}

	return name, keys, nil
}
func (c *DB) Close() {
	c.db.Close()
}
