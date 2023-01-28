package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/Notes-App/operations"
	"github.com/Notes-App/schemas"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CreateNote(c *gin.Context) {
	var input schemas.Note
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	note, err := operations.CreateNote(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, note)
}

func FindNoteByTitle(c *gin.Context) {
	var input schemas.GetNote
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	note, err := operations.FindNoteByTitle(input.Title)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, note)
}

func GetAllNotes(c *gin.Context) {
	notes, err := operations.FindNotes()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, notes)
}

func DeleteNoteByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	n, err := operations.DeleteNote(uuid)
	if err != nil || n != 1 {
		c.JSON(400, err.Error())
		return
	}
	c.Status(200)

}

func UpdateNoteByUUID(c *gin.Context) {
	var input schemas.Note
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	rows, err := operations.UpdateNote(input.UUID, input)
	if err != nil || rows != 1 {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func SignUp(c *gin.Context) {
	var body schemas.User
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	_, err := operations.CreateUser(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	return
}

func Login(c *gin.Context) {
	var body schemas.User
	if c.BindJSON(&body) != nil {
		c.JSON(400, gin.H{"error": "Failed to read body"})
		return
	}
	gotUser, err := operations.FindUser(body)
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
	c.Status(http.StatusOK)
	return
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}
