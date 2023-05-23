package main

import (
	"authentication/controllers"
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"os"
)

func SetupRouter() *gin.Engine {
	ctx := context.Background()
	opt := option.WithCredentialsFile("fiufit-backend-keys.json")
	os.Setenv("FIREBASE_CONFIG", "firebase_config.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln("error initializing app:", err)
	}
	client, err := app.Auth(ctx)
	storage, err := app.Storage(ctx)
	router := gin.Default()
	router.GET("/auth/credentials", controllers.GetCredentials)
	router.GET("/auth/token", controllers.GetToken)
	router.GET("auth/storage/:name", controllers.FileDownload(storage))
	router.POST("/auth", controllers.UserSignUp(client))
	router.POST("auth/storage/:name", controllers.FileUpload(storage))
	router.POST("/auth/tokenLogin", controllers.UserTokenLogin)
	router.POST("/auth/login", controllers.UserLogin)
	router.POST("/auth/loginIDP", controllers.UserVerifyIDPLogin)
	router.POST("/auth/recovery", controllers.PasswordRecovery(client))
	return router
}
