package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mickaelyoshua7674/htmx-study/contact"
)

var cts = contact.ReadJSON()
var maxId = cts.GetMaxId()
var contactErrors = contact.NewContactErrors()

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/contacts")
	})
	router.GET("/contacts", handlerGetContacts)
	router.GET("/contacts/new", handlerNewContact)
	router.POST("/contacts/new", handlerCreateNewContact)
	router.GET("/contacts/:contact_id", handlerShowContact)
	router.GET("/contacts/:contact_id/edit", handlerEditContact)

	router.Run(":8081")
}

func handlerGetContacts(ctx *gin.Context) {
	email := ctx.Request.FormValue("email")
	if email == "" {
		ctx.HTML(http.StatusOK, "index.html", gin.H{"searchEmail":email, "contacts":cts})
	} else {
		ct := cts.HaveEmail(email)
		ctx.HTML(http.StatusOK, "index.html", gin.H{"searchEmail":email, "contacts":ct})
	}
}

func handlerNewContact(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new-contact.html", contact.Contact{})
}

func handlerCreateNewContact(ctx *gin.Context) {
	name := ctx.Request.FormValue("name")
	email := ctx.Request.FormValue("email")
	phone := ctx.Request.FormValue("phone")
	ct := contact.NewContact(maxId+1, name, phone, email, nil)
	maxId++

	cts = append(cts, ct)
	err := cts.WriteJSON()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "new-contact.html", contact.Contact{Errors: contactErrors})
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/contacts")
	}
}

func handlerShowContact(ctx *gin.Context) {
	idString := ctx.Param("contact_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid Id")
		return
	}
	for _, c := range cts {
		if id == c.Id {
			ctx.HTML(http.StatusOK, "show-contact.html", c)
			return
		}
	}
	ctx.String(http.StatusNotFound, "Id not found")
}

func handlerEditContact(ctx *gin.Context) {
	idString := ctx.Param("contact_id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid Id")
		return
	}
	for _, c := range cts {
		if id == c.Id {
			ctx.HTML(http.StatusOK, "edit-contact.html", c)
			return
		}
	}
}