package controller

import (
	"blog-system/config"
	"blog-system/db"
	"blog-system/models"
	"blog-system/pkg/response"
	"blog-system/router/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(g *gin.Context) {
	var loginReq LoginReq
	if err := g.ShouldBindJSON(&loginReq); err != nil {
		config.Logger.Errorf(err.Error())
		response.Fail(g, "参数错误")
		return
	}
	username := loginReq.Username
	pwd := loginReq.Password
	var user models.User
	if err := db.DB.Self.Where("username = ?", username).First(&user).Error; err != nil {
		config.Logger.Errorf(err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(g, "用户名或密码错误")
			return
		}
		response.Fail(g, "系统异常，请稍后再试")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)); err != nil {
		config.Logger.Errorf(err.Error())
		response.Fail(g, "用户名或密码错误")
		return
	}
	tokenString, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		config.Logger.Errorf(err.Error())
		response.Fail(g, "登录失败")
		return
	}
	config.Logger.Info("用户登录", zap.Any("userID", user.ID))
	response.Success(g, tokenString)
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(g *gin.Context) {
	var registerReq RegisterReq
	if err := g.ShouldBindJSON(&registerReq); err != nil {
		config.Logger.Errorf(err.Error())
		response.Fail(g, "注册失败")
		return
	}
	var judgeUser models.User
	db.DB.Self.Where("username=?", registerReq.Username).First(&judgeUser)
	if judgeUser.ID != 0 {
		response.Fail(g, "用户已存在")
		return
	}
	hashedPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if bcryptErr != nil {
		config.Logger.Errorf(bcryptErr.Error())
		response.Fail(g, "注册失败")
		return
	}
	user := models.User{
		Username: registerReq.Username,
		Password: string(hashedPassword),
		Email:    registerReq.Email,
	}
	if err := db.DB.Self.Create(&user); err.Error != nil {
		config.Logger.Errorf(err.Error.Error())
		response.Fail(g, "注册失败")
		return
	}
	response.Success(g, user.ID)
}
