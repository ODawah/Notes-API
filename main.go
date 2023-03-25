package main

import (
	"log"

	"github.com/Notes-App/controllers"
	"github.com/Notes-App/middleware"
	"github.com/Notes-App/models"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var router = gin.Default()
	err := models.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.POST("/create-note", middleware.RequireAuth, controllers.CreateNote)
	router.GET("/find-note", middleware.RequireAuth, controllers.FindNoteByTitle)
	router.GET("/all-notes", middleware.RequireAuth, controllers.GetAllNotes)
	router.DELETE("/delete-note/:uuid", middleware.RequireAuth, controllers.DeleteNoteByUUID)
	router.PATCH("/update-note", middleware.RequireAuth, controllers.UpdateNoteByUUID)

	log.Fatalln(router.Run(":8080"))
}
