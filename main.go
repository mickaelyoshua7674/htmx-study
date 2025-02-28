package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mickaelyoshua7674/htmx-study/contact"
)

var cts = contact.ReadJSON()

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/contacts")
	})
	router.GET("/contacts", handlerGetContacts)
	router.GET("/contacts/new", handlerNewContact)
	router.POST("/contacts/new", handlerCreateNewContact)

	//router.Run(":8081")
}

func handlerGetContacts(ctx *gin.Context) {
	email := ctx.Request.FormValue("email")
	if email == "" {
		ctx.HTML(http.StatusOK, "content", cts)
	} else {
		ct := cts.HaveEmail(email)
		ctx.HTML(http.StatusNotFound, "content", ct)
	}
}

func handlerNewContact(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new-contact", contact.Contact{})
}
func handlerCreateNewContact(ctx *gin.Context) {

}