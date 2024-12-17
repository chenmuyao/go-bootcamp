package web

type ArticleEditReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
