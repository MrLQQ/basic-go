package web

type TopLikeVo struct {
	Id      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	LikeCnt int64  `json:"likeCnt"`
}
