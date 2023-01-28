package main

import (
	"log"

	"github.com/Notes-App/database"
	"github.com/Notes-App/handlers"
	"github.com/Notes-App/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var router = gin.Default()
	database.Connect()

	router.POST("/signup", handlers.SignUp)
	router.POST("/login", handlers.Login)
	router.POST("/create-note", middleware.RequireAuth, handlers.CreateNote)
	router.GET("/find-note", middleware.RequireAuth, handlers.FindNoteByTitle)
	router.GET("/all-notes", middleware.RequireAuth, handlers.GetAllNotes)
	router.DELETE("/delete-note/:uuid", middleware.RequireAuth, handlers.DeleteNoteByUUID)
	router.PATCH("/update-note", middleware.RequireAuth, handlers.UpdateNoteByUUID)
	log.Fatalln(router.Run(":8080"))
}
