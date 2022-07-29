package controllers

import (
	"go-jichu/logic"
	"go-jichu/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	//获取参数以及参数的校验
	p := new(models.Post)
	//从c 上下文中取得当前发起请求的用户id
	userId, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorId = userId

	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindjson err", zap.Any("err", err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//创建贴子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	//获取贴子的id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetPostById(id)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	pageNum, pageSize := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, data)
}

//根据排序返回列表，按创建时间排序，或按分数排序
//获取参数
//去redis中查询id列表
//根据id去数据库查询贴子详细信息
func GetPostListSortHandler(c *gin.Context) {
	//初始化结构体时指定初始参数
	p := &models.ParamPostList{
		PageNum:   1,
		PageLimit: 10,
		Sort:      models.SortTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListSortHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//获取数据
	data, err := logic.GetPostListNew(p)

	if err != nil {
		zap.L().Error("logic.GetPostListNew() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//根据社区去查询帖子列表
// func GetCommunityPostListHandler(c *gin.Context) {
// 	//初始化结构体时指定初始参数
// 	p := &models.ParamCommunityPostList{
// 		ParamPostList: models.ParamPostList{
// 			PageNum:   1,
// 			PageLimit: 10,
// 			Sort:      models.SortTime,
// 		},
// 	}
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("ParamCommunityPostList with invalid params", zap.Error(err))
// 		ResponseError(c, CodeInvalidParam)
// 		return
// 	}
// 	//获取数据
// 	data, err := logic.GetCommunityPostListSort(p)
// 	if err != nil {
// 		zap.L().Error("logic.GetCommunityPostListSort() failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuccess(c, data)
// }
