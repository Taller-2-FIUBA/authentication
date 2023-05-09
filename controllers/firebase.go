package controllers

import (
	"authentication/config"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
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

func sendCustomEmail(email string, username string, link string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "grupocincofiuba.t2@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Recuperacion pass FiuFit")
	m.SetBody("text/html", "Link de recuperación de contraseña para "+username+
		"<br>"+"<a href="+link+">LINK</a>")
	d := gomail.NewDialer("smtp.gmail.com", 587,
		"grupocincofiuba.t2@gmail.com", "hmqfvwlmszqsfhen\n")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return
}

func PasswordRecovery(client *auth.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		email := c.Query("email")
		username := c.Query("username")
		link, err := client.PasswordResetLinkWithSettings(c, email, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Error sending mail"})
			return
		}
		sendCustomEmail(email, username, link)
	}
	return fn
}

func UserSignUp(client *auth.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
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
	posturl := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + config.Cfg.Firebase.Apikey
	body, err := json.Marshal(m)
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Error logging in"})
		return
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil || res.StatusCode != 200 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Error logging in"})
		return
	}
	defer res.Body.Close()
	post := &FireBaseResponse{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		panic(derr)
	}
	c.JSON(http.StatusOK, gin.H{"id": post.LocalId})
}

func UserTokenLogin(c *gin.Context) {
	token := ExtractToken(c)
	if token != "" {
		if !ValidToken(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Expired or invalid token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "No token"})
		return
	}
}
