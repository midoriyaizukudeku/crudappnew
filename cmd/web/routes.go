package main

import (
	"net/http"

	"dream.website/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Notfound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.SessionManager.LoadAndSave, noScrf, app.Authenticated)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.Home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.ViewSnipppet))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protectedRoutes := dynamic.Append(app.requireAuthenticated)

	router.Handler(http.MethodGet, "/snippet/create", protectedRoutes.ThenFunc(app.CreateSnipppet))
	router.Handler(http.MethodPost, "/snippet/create", protectedRoutes.ThenFunc(app.CreateSnipprtPost))
	router.Handler(http.MethodPost, "/user/logout", protectedRoutes.ThenFunc(app.userLogout))
	standard := alice.New(app.recoverPanic, app.logRequest, SecureHeader)

	return standard.Then(router)
}
