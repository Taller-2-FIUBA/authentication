package main

import (
	"authentication/config"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	Data struct {
		ID   string `json:"id"`
		Role string `json:"role"`
	} `json:"data"`
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
	req, _ := http.NewRequest("GET", "/auth/token?role=user&id=pepe", nil)
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
	config.Init()
	router := SetupRouter()
	w := httptest.NewRecorder()
	req1, err := http.NewRequest("GET", "/auth/token?role=user&id=pepe", nil)
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
	assert.Equal(t, "pepe", result.Data.ID)
	assert.Equal(t, "user", result.Data.Role)
}
