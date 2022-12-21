package handlers

import (
	"github.com/Notes-App/operations"
	"github.com/Notes-App/schemas"
	"github.com/gin-gonic/gin"
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
