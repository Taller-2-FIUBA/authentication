package main

import (
	"bytes"
	"context"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	storage2 "firebase.google.com/go/v4/storage"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/mail.v2"
	"io"
	"net/http"
	"time"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UploadedImage struct {
	Image string `json:"image"`
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

type GoogleResponse struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"audience"`
	UserID        string `json:"user_id"`
	Scope         string `json:"scope"`
	ExpiresIn     int32  `json:"expires_in"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	AccessType    string `json:"access_type"`
}

func sendCustomEmail(c *gin.Context, email string, username string, link string) {
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
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Unable to send email"})
		return
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
		sendCustomEmail(c, email, username, link)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Incorrect details for user login"})
		return
	}
	m := userLogin{user.Email, user.Password, true}
	posturl := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + Cfg.Firebase.Apikey
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

func UserVerifyIDPLogin(c *gin.Context) {
	token := ExtractToken(c, false)
	if token != "" {
		url := "https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=" + token
		r, err := http.Get(url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Error communicating with Google"})
			return
		}
		defer r.Body.Close()
		creds := &GoogleResponse{}
		derr := json.NewDecoder(r.Body).Decode(creds)
		if derr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Invalid IDP Token"})
			return
		}
		if !creds.VerifiedEmail || creds.ExpiresIn == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Invalid IDP Token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "No token"})
	return
}

func UserTokenLogin(c *gin.Context) {
	token := ExtractToken(c, true)
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

func FileUpload(storage *storage2.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var image UploadedImage
		if err := c.BindJSON(&image); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Incorrect image encoding"})
			return
		}
		decodedImage, _ := b64.StdEncoding.DecodeString(image.Image)
		fileName := c.Param("name")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
		defer cancel()
		bucket, _ := storage.DefaultBucket()
		object := bucket.Object(fileName)
		wc := object.NewWriter(ctx)
		defer wc.Close()
		if _, err := io.Copy(wc, bytes.NewReader(decodedImage)); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Failed to upload image"})
			return
		}
	}
	return fn
}

func FileDownload(storage *storage2.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		fileName := c.Param("name")
		_, cancel := context.WithTimeout(context.Background(), time.Second*50)
		defer cancel()
		bucket, _ := storage.DefaultBucket()
		object := bucket.Object(fileName)
		rc, err := object.NewReader(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "No image with that name"})
			return
		}
		defer rc.Close()
		data, err := io.ReadAll(rc)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": "Couldn't read object"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"image": data})
	}
	return fn
}
