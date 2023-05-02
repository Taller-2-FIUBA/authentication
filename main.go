package main

import (
	"os"
)

func main() {
	arg := os.Args[1:]
	if len(arg) < 1 {
		print("No port provided, exiting")
		return
	}
	router := SetupRouter()
	err := router.Run(":" + arg[0])
	if err != nil {
		print("Error connecting, exiting")
	}
}
