package main

import (
	"fmt"
	"net/http"

	"github.com/bolatl/lenslocked/controllers"
	"github.com/bolatl/lenslocked/migrations"
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
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	userC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	userC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)
	r.Get("/users/me", userC.CurrentUser)

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}


	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
		// TrustedOrigins takes hosts ("localhost:3000"), not full URLs.
		csrf.TrustedOrigins([]string{"localhost:3000", "127.0.0.1:3000"}),
	)

	fmt.Println("Server starting on :3000...")
	// TODO: Drop plaintextHTTP once we serve over TLS.
	http.ListenAndServe(":3000", plaintextHTTP(csrfMw(umw.SetUser(r))))
}

// plaintextHTTP tells gorilla/csrf that we are serving over cleartext HTTP.
// Without it, csrf assumes an https:// scheme, so a browser's
// "Origin: http://localhost:3000" header never matches and every POST is
// rejected with a 403 "bad origin" error.
func plaintextHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, csrf.PlaintextHTTPRequest(r))
	})
}
