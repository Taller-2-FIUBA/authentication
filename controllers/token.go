package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type header struct {
	Authorization string `header:"Authorization" binding:"required"`
}

var (
	key []byte
	t   *jwt.Token
	s   string
)

type UserClaims struct {
	Role string `json:"role"`
	ID   string `json:"id"`
	jwt.Claims
}

func Encode(c *gin.Context) {
	userId := c.Query("id")
	userRole := c.Query("role")
	key = []byte("secret")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   userId,
			"role": userRole,
		})
	s, _ = t.SignedString(key)
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func GetToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Authorization Header Not Found"})
		return
	}
	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Token contains incorrect data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": claims})
}
