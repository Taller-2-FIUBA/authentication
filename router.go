package main

import (
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
	router.GET("/auth/credentials", GetCredentials)
	router.GET("/auth/token", GetToken)
	router.GET("auth/storage/:name", FileDownload(storage))
	router.POST("/auth", UserSignUp(client))
	router.POST("auth/storage/:name", FileUpload(storage))
	router.POST("/auth/tokenLogin", UserTokenLogin)
	router.POST("/auth/login", UserLogin)
	router.POST("/auth/loginIDP", UserVerifyIDPLogin)
	router.POST("/auth/recovery", PasswordRecovery(client))
	return router
}
