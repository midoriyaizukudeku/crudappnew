package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("err %s\n %s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) Notfound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *Dynamicdata) {
	ts, ok := app.Templatecache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", page))
		return
	}

	buffer := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buffer, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buffer.WriteTo(w)
}

func (app *Application) newTemplateData(r *http.Request) *Dynamicdata {
	_ = r
	return &Dynamicdata{
		CurrentYear:     time.Now().Year(),
		Flash:           app.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *Application) decodePostform(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var Invaliddecodererror *form.InvalidDecoderError

		if errors.As(err, &Invaliddecodererror) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *Application) IsAuthenticated(r *http.Request) bool {
	IsAuthenticated, ok := r.Context().Value(IsAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return IsAuthenticated
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SendVerificationEmail(email, token string) error {
	from := "your-email@example.com"
	password := "your-email-password"
	to := email
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf("Subject: Verify Your Email\n\nClick the link to verify your email: http://localhost:8080/user/verify?token=%s", token))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	return err
}
