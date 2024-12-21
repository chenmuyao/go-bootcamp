package web

type ArticleEditReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticlePublishReq ArticleEditReq

type ArticleWithdrawReq struct {
	ID int64 `json:"id"`
}
