package main

import (
	"log"

	"github.com/Notes-App/database"
	"github.com/Notes-App/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var router = gin.Default()
	database.Connect()

	router.POST("/create-note", handlers.CreateNote)
	router.GET("/find-note", handlers.FindNoteByTitle)
	router.GET("/all-notes", handlers.GetAllNotes)
	router.DELETE("/delete-note/:uuid", handlers.DeleteNoteByUUID)
	router.PATCH("/update-note", handlers.UpdateNoteByUUID)
	log.Fatalln(router.Run(":8080"))
}
