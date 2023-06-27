package main

import (
	"authentication/config"
	"github.com/gin-gonic/gin"
)

func main() {
	if config.Init("config.yml") == false {
		return
	}
	gin.SetMode(gin.ReleaseMode)
	router := SetupRouter()
	err := router.Run(":" + config.Cfg.Server.Port)
	if err != nil {
		print("Error connecting, exiting")
	}
}
