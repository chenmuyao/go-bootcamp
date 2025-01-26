package domain

import "time"

type Article struct {
	ID      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

type Author struct {
	ID   int64
	Name string
}

type ArticleInteractive struct {
	Article Article
	Intr    Interactive
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

const AbstractLen = 128

func (a Article) Abstract() string {
	str := []rune(a.Content)
	if len(str) > AbstractLen {
		str = str[:AbstractLen]
	}
	return string(str)
}
