package controller

import (
	"blog-system/config"
	"blog-system/db"
	"blog-system/models"
	"blog-system/pkg/response"
	"blog-system/router/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	if err := db.DB.Self.Find(&posts).Error; err != nil {
		config.Logger.Error("get all posts error:%v", err)
		response.Fail(c, "获取所有文章失败")
		return
	}
	response.Success(c, posts)
}
func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, "参数有误")
		return
	}
	var posts []models.Post
	if err := db.DB.Self.Find(&posts).Error; err != nil {
		config.Logger.Error("get all posts error:%v", err.Error)
		response.Fail(c, "获取所有文章失败")
		return
	}
	response.Success(c, posts)
}
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		config.Logger.Error("create post error:%v", err.Error)
		response.Fail(c, "参数有误")
		return
	}
	if err := db.DB.Self.Create(&post).Error; err != nil {
		config.Logger.Error("create post error:%v", err)
		response.Fail(c, "文章发发表失败")
		return
	}
	response.Success(c, post)
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		config.Logger.Error("update post error:%v", err.Error)
		response.Fail(c, "更新失败")
		return
	}

	var judgePost models.Post
	if err := db.DB.Self.Where("id = ?", post.ID).First(&judgePost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, "文章不存在")
			return
		}
		config.Logger.Error("UpdatePost error:%v", err)
		response.Fail(c, "获查询文章信息失败")
		return
	}
	currentUserId, err := middleware.GetHeaderUserIdWithError(c)
	if err != nil {
		response.Fail(c, err.Error())
	}
	if currentUserId != judgePost.UserID {
		response.Fail(c, "暂无权限")
	}
	post.CreatedAt = judgePost.CreatedAt
	if err := db.DB.Self.Save(&post).Error; err != nil {
		config.Logger.Error("update post error:%v", err)
		response.Fail(c, "更新失败")
		return
	}
	response.Success(c, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, "参数有误")
		return
	}
	var post models.Post
	if err := db.DB.Self.Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, "文章不存在")
			return
		}
		config.Logger.Error("delete post error:%v", err)
		response.Fail(c, "获查询文章信息失败")
		return
	}
	currentUserId, err := middleware.GetHeaderUserIdWithError(c)
	if err != nil {
		response.Fail(c, err.Error())
	}
	if currentUserId != post.UserID {
		response.Fail(c, "暂无权限")
	}

	if err := db.DB.Self.Delete(&models.Post{}, id).Error; err != nil {
		config.Logger.Error("delete post error:%v", err)
		response.Fail(c, "删除失败")
		return
	}
	response.Success(c, "删除失败")
}
