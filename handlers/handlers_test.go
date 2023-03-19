package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Notes-App/database"
	"github.com/Notes-App/schemas"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestSignUp(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	type test struct {
		name     string
		input    schemas.User
		response string
		code     int
	}

	r := SetUpRouter()
	r.POST("/signup", SignUp)

	tests := []test{
		{name: "normal request", input: schemas.User{Email: "test@test.com", Password: "test123"}, response: `{"message":"signed up"}`, code: http.StatusOK},
		{name: "empty request", input: schemas.User{}, response: `{"error":"Email or Password is empty"}`, code: http.StatusBadRequest},
		{name: "wrong email", input: schemas.User{Email: "test.com", Password: "test123"}, response: `{"error":"invalid email address"}`, code: http.StatusBadRequest},
		{name: "empty password", input: schemas.User{Email: "test@test.com", Password: ""}, response: `{"error":"Email or Password is empty"}`, code: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		jsonValue, _ := json.Marshal(tc.input)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		responseData, _ := ioutil.ReadAll(w.Body)
		assert.Equal(t, tc.code, w.Code)
		assert.Equal(t, tc.response, string(responseData))

	}

}
