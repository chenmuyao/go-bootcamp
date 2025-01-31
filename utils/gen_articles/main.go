package main

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/rand"
)

type Config struct {
	DB *sql.DB
}

type Article struct {
	ID      int64
	Title   string
	Content string

	AuthorID int64
	Status   uint8
	Ctime    int64
	Utime    int64
}

func main() {
	config := &Config{
		DB: initDB(),
	}

	nbGr := 5
	batchSize := 2000

	wg := sync.WaitGroup{}
	wg.Add(nbGr)

	for range nbGr {
		go func() {
			for range batchSize {
				config.generateArticle()
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

func (c *Config) generateArticle() {
	// XXX: Supposing that the DB table is created
	title := randomString(50)
	content := randomString(20000)
	authorID := 1
	status := 2
	ctime := time.Now().UnixMilli()
	utime := time.Now().UnixMilli()

	_, err := c.DB.Exec(
		"INSERT INTO published_articles (title, content, author_id, status, ctime, utime) VALUES (?, ?, ?, ?, ?, ?)",
		title,
		content,
		authorID,
		status,
		ctime,
		utime,
	)
	if err != nil {
		panic(err)
	}
}

func randomString(n int) string {
	// Define the set of characters to choose from.  You can customize this.
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// Create a slice of runes to hold the random string.
	b := make([]rune, n)

	// Seed the random number generator.  This is crucial for getting different
	// random strings each time you call the function.  If you don't seed it,
	// you'll get the same "random" string every time.
	rand.Seed(uint64(time.Now().UnixNano())) // Use current time as seed

	// Generate the random string.
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	// Convert the slice of runes to a string and return it.
	return string(b)
}
