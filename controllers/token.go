package controllers

import (
	"authentication/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	key []byte
	t   *jwt.Token
	s   string
)

type UserClaims struct {
	Role string `json:"role"`
	ID   int    `json:"id"`
	jwt.RegisteredClaims
}

func GetToken(c *gin.Context) {
	userId := c.Query("id")
	numericUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Error converting userID"})
		return
	}
	userRole := c.Query("role")
	if userId == "" || userRole == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Incorrect details"})
		return
	}
	key = []byte("secret")
	claims := UserClaims{
		userRole,
		numericUserId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Cfg.Token.Expiration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ = t.SignedString(key)
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func ValidToken(auth string) bool {
	claims := UserClaims{}
	parsedToken, _ := jwt.ParseWithClaims(auth, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if !parsedToken.Valid {
		return false
	}
	return true
}

func ExtractToken(c *gin.Context, bearer bool) string {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	if bearer {
		splitToken := strings.Split(auth, "Bearer ")
		return splitToken[1]
	}
	return auth
}

func GetCredentials(c *gin.Context) {
	token := ExtractToken(c, true)
	if token == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Authorization Header Not Found"})
		return
	}
	claims := UserClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if !parsedToken.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Invalid token"})
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Message": "Token contains incorrect data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": claims})
}
