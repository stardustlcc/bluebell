package models

//go 内存对齐
type Post struct {
	Id          int64  `json:"id,string" db:"post_id"`
	AuthorId    int64  `json:"author_id" db:"author_id" binding:"required"`
	CommunityId int64  `json:"community_id" db:"community_id" binding:"required"`
	CreateTime  int64  `json:"create_time" db:"create_time"`
	Status      int32  `json:"status" db:"status"`
	Title       string `json:"title" db:"title" binding:"required"`
	Content     string `json:"content" db:"content" binding:"required"`
}

//帖子详情结构体
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name" db:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	*Post                               //嵌入贴子结构体
	*CommunityDetail `json:"community"` //嵌入社区信息
}
