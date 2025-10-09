package main

import (
	"fmt"
	"net/http"

	"github.com/bolatl/lenslocked/controllers"
	"github.com/bolatl/lenslocked/templates"
	"github.com/bolatl/lenslocked/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	userC := controllers.Users{}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)

	fmt.Println("Server starting on :3000...")
	http.ListenAndServe(":3000", r)
}
