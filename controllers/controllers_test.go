package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Notes-App/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/signup", SignUp)
	router.POST("/login", Login)
	router.POST("/create-note", CreateNote)
	router.GET("/find-note", FindNoteByTitle)
	router.GET("/all-notes", GetAllNotes)
	router.DELETE("/delete-note/:uuid", DeleteNoteByUUID)
	router.PATCH("/update-note", UpdateNoteByUUID)

	return router
}

func TestSignUp(t *testing.T) {
	dbErr := models.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer models.CleanUp()
	type test struct {
		name     string
		input    models.User
		response string
		code     int
	}

	r := SetUpRouter()

	tests := []test{
		{name: "normal request", input: models.User{Email: "test@test.com", Password: "test123"}, response: `{"message":"signed up"}`, code: http.StatusOK},
		{name: "empty request", input: models.User{}, response: `{"error":"Email or Password is empty"}`, code: http.StatusBadRequest},
		{name: "wrong email", input: models.User{Email: "test.com", Password: "test123"}, response: `{"error":"invalid email address"}`, code: http.StatusBadRequest},
		{name: "empty password", input: models.User{Email: "test@test.com", Password: ""}, response: `{"error":"Email or Password is empty"}`, code: http.StatusBadRequest},
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

func TestLogin(t *testing.T) {
	dbErr := models.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer models.CleanUp()
	type test struct {
		name     string
		input    models.User
		response string
		code     int
	}

	r := SetUpRouter()
	user := models.User{Email: "test@test.com", Password: "test123"}
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	tests := []test{
		{name: "normal login", input: models.User{Email: "test@test.com", Password: "test123"}, response: `{"message":"logged in"}`, code: http.StatusOK},
		{name: "empty request", input: models.User{}, response: `{"error":"Email or Password is empty"}`, code: http.StatusBadRequest},
		{name: "wrong email", input: models.User{Email: "test.com", Password: "test123"}, response: `{"error":"invalid email address"}`, code: http.StatusBadRequest},
		{name: "empty password", input: models.User{Email: "test@test.com", Password: ""}, response: `{"error":"Email or Password is empty"}`, code: http.StatusBadRequest},
		{name: "wrong user password", input: models.User{Email: "test@test.com", Password: "wrong password"}, response: `{"error":"wrong password"}`, code: http.StatusBadRequest},
		{name: "unregistered user", input: models.User{Email: "test@gdxz.com", Password: "test123"}, response: `{"error":"user not found"}`, code: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		testJsonValue, _ := json.Marshal(tc.input)
		request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testJsonValue))
		writer := httptest.NewRecorder()
		r.ServeHTTP(writer, request)
		responseData, _ := ioutil.ReadAll(writer.Body)
		assert.Equal(t, tc.code, writer.Code)
		assert.Equal(t, tc.response, string(responseData))

	}

}

func TestCreateNote(t *testing.T) {
	dbErr := models.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer models.CleanUp()
	longTitle := "then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	longText := "I'm selfish, I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don' deserve me at my b.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hardimpatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"

	type test struct {
		name  string
		input models.Note
		code  int
	}
	u, _ := models.CreateUser(models.User{Email: "test@test.com", Password: "testpass"})
	r := SetUpRouter()
	tests := []test{
		{name: "normal note input", input: models.Note{Title: "pick up your son", Text: "go to school and pick him up", UserUuid: u.UUID}, code: 200},
		{name: "long note title", input: models.Note{Title: longTitle, Text: "go to school and pick him up"}, code: 400},
		{name: "no note title", input: models.Note{Title: "", Text: "go to school and pick him up"}, code: 400},
		{name: "long note text", input: models.Note{Title: "finish the assignment", Text: longText}, code: 400},
		{name: "no note text", input: models.Note{Title: "cook the dinner", Text: "", UserUuid: u.UUID}, code: 400},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		testJsonValue, _ := json.Marshal(tc.input)
		request, _ := http.NewRequest("POST", "/create-note", bytes.NewBuffer(testJsonValue))
		request.Header.Set("Uuid", u.UUID)
		writer := httptest.NewRecorder()
		r.ServeHTTP(writer, request)
		responseData, _ := ioutil.ReadAll(writer.Body)
		assert.Equal(t, tc.code, writer.Code)
		t.Log(string(responseData))

	}

}

func TestFindNote(t *testing.T) {
	dbErr := models.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer models.CleanUp()
	type test struct {
		name  string
		input models.GetNote
		got   *models.Note
		code  int
	}
	r := SetUpRouter()

	u, _ := models.CreateUser(models.User{Email: "test@test.com", Password: "testpass"})
	n, _ := models.CreateNote(models.Note{Title: "test Title", Text: "test Text", UserUuid: u.UUID})

	tests := []test{
		{name: "normal note title", input: models.GetNote{Title: "test title"}, got: n, code: 201},
		{name: "normal note title with trails", input: models.GetNote{Title: "test title   "}, got: n, code: 201},
		{name: "no note title", input: models.GetNote{Title: ""}, got: nil, code: 400},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		testJsonValue, _ := json.Marshal(tc.input)
		request, _ := http.NewRequest("GET", "/find-note", bytes.NewBuffer(testJsonValue))
		request.Header.Set("Uuid", u.UUID)
		writer := httptest.NewRecorder()
		r.ServeHTTP(writer, request)
		responseData, _ := ioutil.ReadAll(writer.Body)
		var got *models.Note
		err := json.Unmarshal(responseData, &got)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(t, tc.code, writer.Code)
		if tc.got != nil {
			if tc.got.Title != got.Title {
				log.Fatalf("got:%s   expected:%s\n", got.Title, tc.got.Title)
			}
			if tc.got.Text != got.Text {
				log.Fatalf("got:%s   expected:%s\n", got.Text, tc.got.Text)
			}
			if tc.got.UUID != got.UUID {
				log.Fatalf("got:%s   expected:%s\n", got.UUID, tc.got.UUID)
			}
			if tc.got.UserUuid != got.UserUuid {
				log.Fatalf("got:%s   expected:%s\n", got.UserUuid, tc.got.UserUuid)
			}

		}
		t.Log(string(responseData))
	}
}

func TestDeleteNoteByUUID(t *testing.T) {
	dbErr := models.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer models.CleanUp()
	r := SetUpRouter()

	u, _ := models.CreateUser(models.User{Email: "test@test.com", Password: "testpass"})
	n, _ := models.CreateNote(models.Note{Title: "test Title", Text: "test Text", UserUuid: u.UUID})
	url := fmt.Sprintf("/delete-note/%s", n.UUID)
	request, _ := http.NewRequest("DELETE", url, nil)
	request.Header.Set("Uuid", u.UUID)
	writer := httptest.NewRecorder()
	r.ServeHTTP(writer, request)
	if writer.Code != 200 {
		log.Fatalln("error deleting")
	}
	_, err := models.FindNoteByTitle("test Title", u.UUID)
	if err == nil {
		log.Fatalln("note not deleted")
	}
}

func TestUpdateNote(t *testing.T) {
	dbErr := models.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer models.CleanUp()
	r := SetUpRouter()

	u, _ := models.CreateUser(models.User{Email: "test@test.com", Password: "testpass"})
	n, _ := models.CreateNote(models.Note{Title: "test Title", Text: "test Text", UserUuid: u.UUID})
	n.Title = "updated title"
	n.Text = "updated text"
	testBody, _ := json.Marshal(n)
	request, _ := http.NewRequest("PATCH", "/update-note", bytes.NewBuffer(testBody))
	request.Header.Set("Uuid", u.UUID)
	writer := httptest.NewRecorder()
	r.ServeHTTP(writer, request)

	got, _ := models.FindNoteByTitle("updated title", u.UUID)
	if got.Text != n.Text {
		log.Fatal("un updated text")
	}
	if got.Title != n.Title {
		log.Fatal("un updated title")
	}
}
