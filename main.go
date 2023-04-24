package main

import (
	"authentication/controllers"
	"context"
	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
)

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("fiufit-backend-keys.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln("error initializing app:", err)
	}
	client, err := app.Auth(ctx)
	router := gin.Default()
	router.GET("/auth/credentials", controllers.GetToken)
	router.GET("/auth/token", controllers.Encode)
	router.POST("/auth", controllers.UserSignUp(client))
	router.POST("/auth/login", controllers.UserLogin(client))
	router.Run("localhost:8082")
}
