package main

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Response struct {
	Data struct {
		ID   int    `json:"id"`
		Role string `json:"role"`
	} `json:"data"`
}

type Image struct {
	Image string `json:"image"`
}

func TestCantGetTokenWithNoDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/token", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func TestCantGetTokenWithIncorrectDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/token?wrong=pepe", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func TestCanGetTokenWithProperDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/token?role=user&id=1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}

func TestCantGetCredentialsWithInvalidToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/credentials", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func TestCanGetCredentialsWithValidToken(t *testing.T) {
	Init("config.yml")
	router := SetupRouter()
	w := httptest.NewRecorder()
	req1, err := http.NewRequest("GET", "/auth/token?role=user&id=1", nil)
	router.ServeHTTP(w, req1)
	var j map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&j)
	if err != nil {
		panic(err)
	}
	req2, _ := http.NewRequest("GET", "/auth/credentials", nil)
	req2.Header.Set("Authorization", "Bearer "+j["data"].(string))
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 1, result.Data.ID)
	assert.Equal(t, "user", result.Data.Role)
}

func TestCanLoginWithValidToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/auth/token?role=user&id=1", nil)
	router.ServeHTTP(w, req1)
	var j map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&j)
	if err != nil {
		panic(err)
	}
	req2, _ := http.NewRequest("POST", "/auth/tokenLogin", nil)
	req2.Header.Set("Authorization", "Bearer "+j["data"].(string))
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 200, w.Code)
}

func TestCantLoginWithInvalidToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	key := []byte("secret")
	claims := UserClaims{
		"user",
		1,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	finishedToken, _ := token.SignedString(key)
	req2, _ := http.NewRequest("POST", "/auth/tokenLogin", nil)
	req2.Header.Set("Authorization", "Bearer "+finishedToken)
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"Message\":\"Expired or invalid token\"}", w.Body.String())
}

func TestCantLoginWithTokenIfTheresNoToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth/tokenLogin", nil)
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"No token\"}", w.Body.String())
}

func TestCantSendRecoveryEmailWithNoDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth/recovery", nil)
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "{\"Message\":\"Error sending mail\"}", w.Body.String())
}

func TestCantLoginThroughIDPWithNoToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth/loginIDP", nil)
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"Message\":\"No token\"}", w.Body.String())
}

func TestCantLoginThroughIDPWithInvalidToken(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth/loginIDP", nil)
	req2.Header.Set("Authorization", "Bearer "+"invalid token")
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "{\"Message\":\"Invalid IDP Token\"}", w.Body.String())
}

func TestCantLoginWithoutUserDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth/login", nil)
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"Incorrect details for user login\"}", w.Body.String())
}

func TestCantLoginWithNonExistentUser(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	body := User{
		"nosuchemail",
		"fakepassword",
	}
	reqBody, _ := json.Marshal(body)
	req2, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"Error logging in\"}", w.Body.String())
}

func TestCantCreateUserWithoutUserDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/auth", nil)
	router.ServeHTTP(w, req2)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"Incorrect details for user creation\"}", w.Body.String())
}

func TestCantCreateUserWithInvalidUserDetails(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	body := User{
		"nosuchemail",
		"fakepassword",
	}
	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"Error creating user\"}", w.Body.String())
}

func TestCantDownloadFileWithIncorrectName(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/storage/non_existent_file ", nil)
	router.ServeHTTP(w, req)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"No image with that name\"}", w.Body.String())
}

func TestCantUploeadImageWithWrongFormat(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/storage/testimage ", nil)
	router.ServeHTTP(w, req)
	var result Response
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"Message\":\"Incorrect image encoding\"}", w.Body.String())
}

func TestCanUploadAndDownloadImage(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	body := Image{
		"testimage",
	}
	reqBody, _ := json.Marshal(body)
	req1, _ := http.NewRequest("POST", "/auth/storage/testimage", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req1)
	var result Response
	req2, _ := http.NewRequest("GET", "/auth/storage/testimage", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req2)
	json.Unmarshal([]byte(w.Body.String()), &result)
	assert.Equal(t, 200, w.Code)
}

func TestCantOpenNonExistentFile(t *testing.T) {
	assert.Equal(t, Init("fake.yml"), false)
}

func TestCanOpenPresentFile(t *testing.T) {
	assert.Equal(t, Init("config.yml"), true)
}
