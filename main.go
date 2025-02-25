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
	router.GET("/contacts")

	//router.Run(":8081")


}

func GetContacts(ctx *gin.Context) {
	email := ctx.Request.FormValue("email")
	if email == "" {
	}
	ct := cts.HaveEmail(email)

	ctx.HTML(http.StatusOK, "content", ct)
}