package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLogin struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type FireBaseResponse struct {
	IdToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalId      string `json:"localId"`
	Registered   bool   `json:"registered"`
}

func UserSignUp(client *auth.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			print(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Incorrect details for user creation"})
			return
		}
		params := (&auth.UserToCreate{}).
			Email(user.Email).
			Password(user.Password)
		newUser, err := client.CreateUser(context.Background(), params)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Error creating user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": newUser.UID})
	}
	return fn
}

func UserLogin(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		print(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Incorrect details for user login"})
		return
	}
	m := userLogin{user.Email, user.Password, true}
	posturl := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyDIuhP0PN2JyRRAdqegXBzm_YO-HKPgjaQ"
	body, err := json.Marshal(m)
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil || res.StatusCode != 200 {
		panic(err)
	}
	defer res.Body.Close()
	post := &FireBaseResponse{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		panic(derr)
	}
	c.JSON(http.StatusOK, gin.H{"id": post.IdToken})
}
