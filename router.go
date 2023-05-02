package main

import (
	"authentication/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//ctx := context.Background()
	//opt := option.WithCredentialsFile("fiufit-backend-keys.json")
	//app, err := firebase.NewApp(ctx, nil, opt)
	//if err != nil {
	//log.Fatalln("error initializing app:", err)
	//}
	//client, err := app.Auth(ctx)
	router := gin.Default()
	router.GET("/auth/credentials", controllers.GetCredentials)
	router.GET("/auth/token", controllers.GetToken)
	//router.POST("/auth", controllers.UserSignUp(client))
	//router.POST("/auth/login", controllers.UserLogin)
	return router
}
