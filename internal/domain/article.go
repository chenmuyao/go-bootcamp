package domain

type Article struct {
	ID      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type Author struct {
	ID   int64
	Name string
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)
