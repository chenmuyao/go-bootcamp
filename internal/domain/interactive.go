package domain

type Interactive struct {
	Biz        string
	BizID      int64
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Liked      bool
	Collected  bool
}
