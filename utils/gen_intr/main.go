package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/exp/rand"
)

type Config struct {
	DB *sql.DB
}

type Interactive struct {
	ID int64

	// <biz_id, biz>
	BizID      int64
	Biz        string
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Utime      int64
	Ctime      int64
}

func main() {
	config := &Config{
		DB: initDB(),
	}

	batchSize := 2000

	for i := range batchSize {
		config.generateArticle(i)
	}
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:13316)/wetravel")
	if err != nil {
		panic(err)
	}

	return db
}

func (c *Config) generateArticle(i int) {
	// XXX: Supposing that the DB table is created
	bizID := i
	biz := "article"
	readCnt := rand.Int63()
	likeCnt := rand.Int63()
	collectCnt := rand.Int63()
	ctime := time.Now().UnixMilli()
	utime := time.Now().UnixMilli()

	_, err := c.DB.Exec(
		"INSERT INTO interactives (biz_id, biz, read_cnt, like_cnt, collect_cnt, ctime, utime) VALUES (?, ?, ?, ?, ?, ?, ?)",
		bizID,
		biz,
		readCnt,
		likeCnt,
		collectCnt,
		ctime,
		utime,
	)
	if err != nil {
		panic(err)
	}
}
