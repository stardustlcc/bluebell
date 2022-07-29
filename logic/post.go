package logic

import (
	"fmt"
	"go-jichu/dao/mysql"
	"go-jichu/dao/redis"
	"go-jichu/models"
	"go-jichu/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成post_id
	p.Id = int64(snowflake.GetID())
	//保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.Id, p.CommunityId)
	return
}

// func GetPostById(id int64) (data *models.Post, err error) {
// 	return mysql.GetPostById(id)
// }

func GetPostById(id int64) (data *models.ApiPostDetail, err error) {
	//查询数据 组合接口想要的数据

	post, err := mysql.GetPostById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostById(id) failed", zap.Int64("pid", id), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
		return
	}
	//根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		zap.L().Error(" mysql.GetCommunityDetailById(post.CommunityId) failed", zap.Int64("communityId", post.CommunityId), zap.Error(err))
		return
	}
	//拼装数据
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(pageNum, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error(" mysql.GetCommunityDetailById(post.CommunityId) failed", zap.Int64("communityId", post.CommunityId), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostListSort(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIdsInSort(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostIdsInSort return 0 len data")
		return
	}

	//返回的数据要按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	//提前查询好每篇贴子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将贴子的作者以及分区信息查询出来填充到贴子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error(" mysql.GetCommunityDetailById(post.CommunityId) failed", zap.Int64("communityId", post.CommunityId), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	fmt.Println(data)
	return
}

func GetCommunityPostListSort(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	ids, err := redis.GetCommunityPostIdsInSort(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostIdsInSort return 0 len data")
		return
	}

	//返回的数据要按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	//提前查询好每篇贴子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//将贴子的作者以及分区信息查询出来填充到贴子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId) failed", zap.Int64("author_id", post.AuthorId), zap.Error(err))
			continue
		}
		//根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error(" mysql.GetCommunityDetailById(post.CommunityId) failed", zap.Int64("communityId", post.CommunityId), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	fmt.Println(data)
	return
}

//根据请求参数的不同，执行不同的逻辑
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityId == 0 {
		//查所有
		data, err = GetPostListSort(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostListSort(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
