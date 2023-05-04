package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

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

func GetToken(c *gin.Context) {
	print(c.Request)
	userId := c.Query("id")
	userRole := c.Query("role")
	if userId == "" || userRole == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Incorrect details"})
	}
	key = []byte("secret")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   userId,
			"role": userRole,
		})
	s, _ = t.SignedString(key)
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func GetCredentials(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Authorization Header Not Found"})
	}
	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Token contains incorrect data"})
	}
	c.JSON(http.StatusOK, gin.H{"data": claims})
}
