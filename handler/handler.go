package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/mickaelyoshua7674/htmx-study/contact"
	"github.com/mickaelyoshua7674/htmx-study/view"
)

func Render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}
func HandleErrorRender(err error) {
	if err != nil {
		log.Fatalf("Could not render :%v", err)
	}
}

func GetContacts(ctx *gin.Context) {
	time.Sleep(2*time.Second)
	cts := contact.ReadJSON()

	var page int
	var err error
	strPage := ctx.Request.FormValue("page")
	if strPage != "" {
		page, err = strconv.Atoi(strPage)
		if err != nil || page <= 0 {
			ctx.String(http.StatusBadRequest, "Invalid page")
			return
		}
	} else {
		page = 1
	}

	query := ctx.Request.FormValue("query")
	cts = cts.GetByQuery(query)
	if ctx.GetHeader("HX-Trigger") == "search" {
		err = Render(ctx, http.StatusOK, view.IndexTr(cts))
		HandleErrorRender(err)
	} else {
		err = Render(ctx, http.StatusOK, view.Index(query, cts, page))
		HandleErrorRender(err)
	}
}

func GetCount(ctx *gin.Context) {
	cts := contact.ReadJSON()
	time.Sleep(2*time.Second)
	ctx.String(http.StatusOK, cts.GetCountStr() + " total contacts")
}

func FormNewContact(ctx *gin.Context) {
	err := Render(ctx, http.StatusOK, view.NewContact(contact.Contact{}))
	HandleErrorRender(err)
}

func CreateNewContact(ctx *gin.Context) {
	cts := contact.ReadJSON()
	maxId := cts.GetMaxId()
	maxId++

	name := ctx.Request.FormValue("name")
	email := ctx.Request.FormValue("email")
	phone := ctx.Request.FormValue("phone")
	ct := contact.NewContact(maxId+1, name, phone, email)

	cts = append(cts, ct)
	err := cts.WriteJSON()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error saving contacts: %v", err)
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

	ct := cts.GetContactById(id)
	if ct.Email == "" {
		ctx.String(http.StatusNotFound, "Id not found")
	} else {
		err = Render(ctx, http.StatusOK, view.ShowContact(ct))
		HandleErrorRender(err)
	}
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
		err := Render(ctx, http.StatusOK, view.EditContact(cts[index]))
		HandleErrorRender(err)
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