package main

import (
	"blog-system/config"
	"blog-system/db"
	"blog-system/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
	pflag.Parse()

	// 初始化配置
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}
	g := gin.Default()
	middlewares := []gin.HandlerFunc{}

	// routers
	router.Load(g, middlewares...)
	// 初始化数据库
	db.InitDB()
	gin.SetMode(viper.GetString("runmode"))

	config.Logger.Infof("Start to requests on http address: %s", viper.GetString("addr"))
	config.Logger.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}
