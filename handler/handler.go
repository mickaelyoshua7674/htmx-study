package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mickaelyoshua7674/htmx-study/contact"
)

var contactErrors = contact.NewContactErrors()

func GetContacts(ctx *gin.Context) {
	cts := contact.ReadJSON()

	email := ctx.Request.FormValue("email")
	if email == "" {
		ctx.HTML(http.StatusOK, "index.html", gin.H{"searchEmail": email, "contacts": cts})
	} else {
		ct := cts.HaveEmail(email)
		ctx.HTML(http.StatusOK, "index.html", gin.H{"searchEmail": email, "contacts": ct})
	}
}

func NewContact(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "new-contact.html", contact.Contact{})
}

func CreateNewContact(ctx *gin.Context) {
	cts := contact.ReadJSON()
	maxId := cts.GetMaxId()

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

func ShowContact(ctx *gin.Context) {
	cts := contact.ReadJSON()

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

func FormEditContact(ctx *gin.Context) {
	cts := contact.ReadJSON()

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
	ctx.String(http.StatusNotFound, "Id not found")
}

func EditContact(ctx *gin.Context) {
	cts := contact.ReadJSON()

	name := ctx.Request.FormValue("name")
	phone := ctx.Request.FormValue("phone")
	email := ctx.Request.FormValue("email")

	idString := ctx.Param("contact_id")
	requestId, err := strconv.Atoi(idString)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid Id")
		return
	}

	index := cts.GetIndexById(requestId)
	if index == -1 {
		ctx.String(http.StatusNotFound, "Id not found")
		return
	}

	cts[index].Update(name, phone, email)

	err = cts.WriteJSON()
	if err != nil {
		cts[index].Errors["email"] = err
		ctx.HTML(http.StatusInternalServerError, "edit-contact.html", cts[index])
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/contacts/"+idString)
	}
}

func DeleteContact(ctx *gin.Context) {
	cts := contact.ReadJSON()

	idString := ctx.Param("contact_id")
	id, _ := strconv.Atoi(idString)
	// No need to verify the id errors because
	// it was already verifyed in the edit handler functions

	cts.DeleteById(id)

	err := cts.WriteJSON()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error saving changes")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/contacts")
	}
}