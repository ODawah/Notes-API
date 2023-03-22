package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Notes-App/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CreateNote(c *gin.Context) {
	uuid, err := c.Cookie("uuid")
	if err != nil {
		c.JSON(401, gin.H{"error": "UnAuthorized"})
		return
	}
	var input models.Note
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	input.UserUuid = uuid
	note, err := models.CreateNote(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, note)
}

func FindNoteByTitle(c *gin.Context) {
	uuid, err := c.Cookie("uuid")
	if err != nil {
		c.JSON(401, gin.H{"error": "UnAuthorized"})
		return
	}
	var input models.GetNote
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	note, err := models.FindNoteByTitle(input.Title, uuid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, note)
}

func GetAllNotes(c *gin.Context) {
	uuid, err := c.Cookie("uuid")
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	notes, err := models.FindNotes(uuid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, notes)
}

func DeleteNoteByUUID(c *gin.Context) {
	userUuid, err := c.Cookie("uuid")
	if err != nil {
		c.JSON(401, gin.H{"error": "UnAuthorized"})
		return
	}
	uuid := c.Param("uuid")
	err = models.DeleteNote(uuid, userUuid)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.Status(200)

}

func UpdateNoteByUUID(c *gin.Context) {
	var input models.Note
	userUuid, err := c.Cookie("uuid")
	if err != nil {
		c.JSON(401, gin.H{"error": "UnAuthorized"})
		return
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	input.UserUuid = userUuid
	err = models.UpdateNote(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func SignUp(c *gin.Context) {
	var body models.User
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is empty"})
		return
	}
	_, err := models.CreateUser(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed up"})
	return
}

func Login(c *gin.Context) {
	var body models.User
	if c.BindJSON(&body) != nil {
		c.JSON(400, gin.H{"error": "Failed to read body"})
		return
	}
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is empty"})
		return
	}
	gotUser, err := models.FindUser(body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": gotUser.UUID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.SetCookie("uuid", gotUser.UUID, 3600, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
	return
}
