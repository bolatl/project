package controllers

import (
	"fmt"
	"net/http"

	"github.com/bolatl/lenslocked/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	SessionService *models.SessionService
	UserService    *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println("123")
	user, err := u.UserService.Create(email, password)
	fmt.Println("123")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in creating a user", http.StatusInternalServerError)
		return
	}
	fmt.Println("123")
	session, err := u.SessionService.Create(user.ID)
	fmt.Println("123")
	if err != nil {
		fmt.Println(err)
		// long term here warning to user should be sent
		http.Redirect(w, r, "/signin", http.StatusAccepted)
		return
	}
	fmt.Println("123")
	setCookie(w, CookieSession, session.Token)
	fmt.Println("123")
	http.Redirect(w, r, "/users/me", http.StatusFound)
	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in logging in", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// long term here warning to user should be sent
		http.Redirect(w, r, "/signin", http.StatusAccepted)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	fmt.Fprintf(w, "Successful login: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "curr user: %s\n", user)
	fmt.Fprintf(w, "Headers: %+v\n", r.Header)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong, try again.", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
