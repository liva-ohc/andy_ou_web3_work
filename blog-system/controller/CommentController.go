package controller

import (
	"blog-system/config"
	"blog-system/db"
	"blog-system/models"
	"blog-system/pkg/response"
	"github.com/gin-gonic/gin"
)

func GetAllCommentByPostId(c *gin.Context) {
	postId := c.Param("postId")
	if postId == "" {
		response.Fail(c, "参数有误")
		return
	}

	var comments []models.Comment
	if err := db.DB.Self.Where("postId=?", postId).Find(&comments); err != nil {
		response.Fail(c, "系统异常，请稍后再试")
		return
	}

	response.Success(c, comments)

}

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		config.Logger.Error("create comment error:%v", err)
		response.Fail(c, "评论发表失败")
		return
	}
	if err := db.DB.Self.Create(&comment).Error; err != nil {
		config.Logger.Error("create comment error:%v", err)
		response.Fail(c, "评论发表失败")
		return
	}
	response.Success(c, comment.ID)
}
