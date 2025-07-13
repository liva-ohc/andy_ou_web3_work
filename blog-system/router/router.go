package router

import (
	"blog-system/controller"
	"blog-system/router/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(mw...)

	g.POST("/login", controller.Login)
	g.POST("/register", controller.Register)

	// 文章路由组
	postGroup := g.Group("/posts")
	{
		// 获取所有文章 (无需认证)
		postGroup.GET("", controller.GetAllPosts)

		// 获取指定ID的文章 (无需认证)
		postGroup.GET("/:id", controller.GetPostByID)

		// 创建文章 (需要认证)
		postGroup.POST("", middleware.JwtAuthMiddleware(), controller.CreatePost)

		// 更新文章 (需要认证)
		postGroup.PUT("/:id", middleware.JwtAuthMiddleware(), controller.UpdatePost)

		// 删除文章 (需要认证)
		postGroup.DELETE("/:id", middleware.JwtAuthMiddleware(), controller.DeletePost)
	}

	commentGroup := g.Group("comment")
	{
		//根据id获取文章
		commentGroup.GET("/:id", controller.GetAllCommentByPostId)
		//发表评论
		commentGroup.POST("", middleware.JwtAuthMiddleware(), controller.CreateComment)
	}
	pprof.Register(g)
	return g
}
