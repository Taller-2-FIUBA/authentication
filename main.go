package main

import (
	"authentication/config"
)

func main() {
	if config.Init() == false {
		return
	}
	//gin.SetMode(gin.ReleaseMode)
	router := SetupRouter()
	err := router.Run(":" + config.Cfg.Server.Port)
	if err != nil {
		print("Error connecting, exiting")
	}
}
