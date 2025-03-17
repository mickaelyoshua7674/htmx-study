package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mickaelyoshua7674/htmx-study/contact"
	"github.com/mickaelyoshua7674/htmx-study/view"
)

func GetContacts(ctx *gin.Context) {
	cts := contact.ReadJSON()
	view.Index(cts)

	email := ctx.Request.FormValue("email")
	id := cts.GetIdByEmail(email)

	if email == "" {
		ctx.HTML(http.StatusOK, "index.html", gin.H{"searchEmail": email, "contacts": cts})
	} else {
		if id == -1 {
			ctx.String(http.StatusNotFound, "Contact not found")
		} else {
			ct := cts.GetContactById(id)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"searchEmail": email, "contacts": ct})
		}
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
	ct := contact.NewContact(maxId+1, name, phone, email)
	maxId++

	cts = append(cts, ct)
	err := cts.WriteJSON()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "new-contact.html", contact.Contact{})
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

	index := cts.GetIndexById(id)
	if index != -1 {
		ctx.HTML(http.StatusOK, "edit-contact.html", cts[index])
	} else {
		ctx.String(http.StatusNotFound, "Id not found")
	}
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
		// Using "See Other" (303) because by default 301 and 302 will send
		// the same request method (in this case DELETE) when redirecting.
		// Since is needed a GET request to "/contacts" to redirect correctly
		// will be send the status "See Other".
		ctx.Redirect(http.StatusSeeOther, "/contacts")
	}
}

func ValidateEmail(ctx *gin.Context) {
	cts := contact.ReadJSON()

	idString := ctx.Param("contact_id")
	id, _ := strconv.Atoi(idString)
	// No need to verify the id errors because
	// it was already verifyed in the edit handler functions
	c := cts.GetContactById(id)

	requestEmail := ctx.Request.FormValue("email")
	searchId := cts.GetIdByEmail(requestEmail)
	if searchId != -1 && requestEmail != c.Email {
		ctx.String(http.StatusBadRequest, "Email must be Unique")
	} else {
		ctx.String(http.StatusOK, "")
	}
}