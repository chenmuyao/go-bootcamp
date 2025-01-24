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

type ArticleVO struct {
	ID         int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	Content    string `json:"content,omitempty"`
	AuthorID   int64  `json:"authorId,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`

	ReadCnt    int64 `json:"readCnt,omitempty"`
	LikeCnt    int64 `json:"likeCnt,omitempty"`
	CollectCnt int64 `json:"collectCnt,omitempty"`
	Liked      bool  `json:"liked"`
	Collected  bool  `json:"collected"`
}

type Like struct {
	ID   int64 `json:"id"`
	Like bool  `json:"liked"`
}

type Collect struct {
	ID int64 `json:"id"`
	// collection id
	CID       int64 `json:"cid"`
	Collected bool  `json:"collected"`
}
