package main

import (
	"fmt"
	"net/http"

	"github.com/bolatl/lenslocked/controllers"
	"github.com/bolatl/lenslocked/models"
	"github.com/bolatl/lenslocked/templates"
	"github.com/bolatl/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	userC := controllers.Users{
		UserService: &userService,
	}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	userC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.Get("/users/me", userC.CurrentUser)

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"https://127.0.0.1:3000",
			"https://localhost:3000",
		}),
	)
	fmt.Println("Server starting on :3000...")
	http.ListenAndServe(":3000", csrfMw(r))
}
