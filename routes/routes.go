package routes

import (
	"go-jichu/controllers"
	"go-jichu/logger"
	"go-jichu/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册
	v1.POST("/register", controllers.RegisterHandler)
	//登录
	v1.POST("/login", controllers.LoginHandler)
	//jwt认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		//专题
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		//贴子
		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/post-list", controllers.GetPostListHandler)

		//投票
		v1.POST("/vote", controllers.PostVoteController)
		//根据时间或分数获取帖子列表
		v1.GET("/post-list-sort", controllers.GetPostListSortHandler)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
