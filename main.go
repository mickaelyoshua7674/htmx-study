package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mickaelyoshua7674/htmx-study/handler"
)

func main() {
	router := gin.Default()
	router.Static("/static", "./static")

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/contacts")
	})
	router.GET("/contacts", handler.GetContacts)
	router.GET("/contacts/new", handler.FormNewContact)
	router.POST("/contacts/new", handler.CreateNewContact)
	router.GET("/contacts/:contact_id", handler.ShowContact)
	router.GET("/contacts/:contact_id/edit", handler.FormEditContact)
	router.POST("/contacts/:contact_id/edit", handler.EditContact)
	router.GET("/contacts/:contact_id/email", handler.ValidateEmail)
	router.DELETE("/contacts/:contact_id", handler.DeleteContact)

	router.Run(":8081")
}