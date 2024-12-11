package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Config struct {
	DB *sql.DB
}

type User struct {
	ID       int64
	Email    string
	Password string

	Ctime int64
	Utime int64

	Name     string
	Birthday int64
	Profile  string
}

func main() {
	config := &Config{
		DB: initDB(),
	}

	nbGr := 100
	batchSize := 100000

	wg := sync.WaitGroup{}
	wg.Add(nbGr)

	for range nbGr {
		go func() {
			for range batchSize {
				id, _ := uuid.NewRandom()
				user := User{
					Email:    fmt.Sprintf("%s@test.com", id.String()),
					Password: "password123!",
				}
				err := config.writeUser(&user)
				if err != nil {
					fmt.Println(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:13316)/wetravel")
	if err != nil {
		panic(err)
	}

	return db
}

func (c *Config) writeUser(user *User) error {
	// XXX: Supposing that the DB table is created
	_, err := c.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)",
		user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
