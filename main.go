package main

import (
	"electronic-gallery/global"
	"electronic-gallery/internal/cache"
	"electronic-gallery/internal/model"
	"electronic-gallery/internal/routers"
	"electronic-gallery/pkg/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	err := setupSetting()
	if err != nil {
		panic(err)
	}
	err = model.SetupDBEngine()
	if err != nil {
		panic(err)
	}
	err = cache.SetupRedis()
	if err != nil {
		panic(err)
	}
}

// @title 电子相册系统
// @version 1.0.0
// @description
// @termsOfService https://github.com/nanjingblue/electronic-gallery-backend
// @host 127.0.0.1:8080
// @BasePath /api/v1
func main() {
	//fmt.Println("hello world")
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:        ":" + global.ServerSetting.HttpPort,
		Handler:     router,
		ReadTimeout: global.ServerSetting.ReadTimeOut,
		//WriteTimeout:   global.ServerSetting.WriteTimeOut,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		return
	}
	return
}

func setupSetting() interface{} {
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Jwt", &global.JwtSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("OSS", &global.OSSSetting)
	if err != nil {
		return err
	}
	return nil
}
