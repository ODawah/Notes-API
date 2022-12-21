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

// err := database.Connect()
// if err != nil {
// fmt.Println(err)
// return
// }
// note, err := operations.CreateNote(schemas.Note{Title: "hello", Text: "world"})
// if err != nil {
// fmt.Println(err)
// return
// }
// fmt.Println(note)
// note, err = operations.CreateNote(schemas.Note{Title: "hellsadfsfsdfo", Text: "woasdfasdfgdfbharld"})
// if err != nil {
// fmt.Println(err)
// }
// notes, err := operations.FindNotes()
// fmt.Println(notes)
