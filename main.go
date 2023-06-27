package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	if Init("config.yml") == false {
		return
	}
	gin.SetMode(gin.ReleaseMode)
	router := SetupRouter()
	err := router.Run(":" + Cfg.Server.Port)
	if err != nil {
		print("Error connecting, exiting")
	}
}
