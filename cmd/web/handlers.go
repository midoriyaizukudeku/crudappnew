package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"dream.website/internal/model"
	"dream.website/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type UserSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
type userloginform struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "index.html", data)
	for name := range app.Templatecache {
		app.infoLog.Printf("Template cached: %s", name)
	}
}

func (app *Application) ViewSnipppet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.Notfound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, model.ErrnoRecord) {
			app.Notfound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.SingleSnippet = snippet
	app.render(w, http.StatusOK, "view.html", data)
}

func (app *Application) CreateSnipprtPost(w http.ResponseWriter, r *http.Request) {

	var form SnippetCreateForm

	err := app.decodePostform(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "this firld cant be empty")
	form.CheckField(validator.MaxChar(form.Title, 100), "title", "cant be more then 100 charecters")
	form.CheckField(validator.NotBlank(form.Content), "content", "THis foeld cant be empty")
	form.CheckField(validator.PermittedValues(form.Expires, 7, 30, 365), "expires", "this field must be 7 , 30, 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "snippet created successfully")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *Application) CreateSnipppet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = SnippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &UserSignUpForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	var form UserSignUpForm
	err := app.decodePostform(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "name cant be empty")
	form.CheckField(validator.NotBlank(form.Email), "email", "this field cant be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "this field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "this field cant be empty")
	form.CheckField(validator.MinChar(form.Password, 8), "password", "minimum 8 charecters required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, model.DuplicateEmail) {
			form.AddFieldError("email", "Same email already in use ")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.SessionManager.Put(r.Context(), "flash", "Your signup was successful please log in ")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}
func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userloginform{}
	app.render(w, http.StatusSeeOther, "login.html", data)
}
func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userloginform

	err := app.decodePostform(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "this field can't be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "not valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "this field csnt be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusSeeOther, "login.html", data)
		return
	}
	id, err := app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, model.InvalidCredientials) {
			form.NonFieldError = append(form.NonFieldError, "email or pass incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusSeeOther, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "authenticatedUserId", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

}
func (app *Application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.SessionManager.Remove(r.Context(), "authenticatedUserId")
	app.SessionManager.Put(r.Context(), "flash", "You have been logget out ")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
