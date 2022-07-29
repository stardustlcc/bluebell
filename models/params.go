package models

const (
	SortTime  = "time"
	SortScore = "score"
)

//定义请求的参数结构体

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	PostId    string `json:"post_id" binding:"required"`
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` //赞成票1 反对票-1 取消投票0  注意：required 0 值会忽略
}

//paramPostList 获取贴子列表query string 参数
type ParamPostList struct {
	CommunityId int64  `json:"community_id" form:"community_id"` //可以为空
	PageNum     int64  `json:"page_num" form:"page_num"`
	PageLimit   int64  `json:"page_limit" form:"page_limit"`
	Sort        string `json:"sort" form:"sort"`
}

// type ParamCommunityPostList struct {
// 	ParamPostList

// }
