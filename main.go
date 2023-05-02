package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	arg := os.Args[1:]
	if len(arg) < 1 {
		print("No port provided, exiting")
		return
	}
	gin.SetMode(gin.ReleaseMode)
	router := SetupRouter()
	err := router.Run(":" + arg[0])
	if err != nil {
		print("Error connecting, exiting")
	}
}
