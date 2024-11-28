package domain

type Interactive struct {
	ReadCnt    int64
	BizId      int64
	Biz        string
	LikeCnt    int64
	CollectCnt int64
	Liked      bool
	Collected  bool
}
