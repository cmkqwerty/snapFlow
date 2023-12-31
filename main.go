package main

import (
	"fmt"
	"github.com/cmkqwerty/snapFlow/configs"
	"github.com/cmkqwerty/snapFlow/controllers"
	"github.com/cmkqwerty/snapFlow/migrations"
	"github.com/cmkqwerty/snapFlow/models"
	"github.com/cmkqwerty/snapFlow/templates"
	"github.com/cmkqwerty/snapFlow/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"strconv"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dbCfg := struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		SSLMode  string
	}{
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBSSLMode,
	}

	db, err := models.Open(dbCfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	smtpCfgPort, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		log.Fatal("SMTP Config Port: %w", err)
	}

	smtpCfg := struct {
		Host     string
		Port     int
		Username string
		Password string
	}{
		config.SMTPHost,
		smtpCfgPort,
		config.SMTPUsername,
		config.SMTPPassword,
	}

	serverAddress := config.ServerAddress

	csrfKey := config.CSRFKey
	csrfSecure, err := strconv.ParseBool(config.CSRFSecure)
	if err != nil {
		log.Fatal("csrfSecure param: %w", err)
	}

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	passwordResetService := &models.PasswordResetService{
		DB: db,
	}

	galleryService := &models.GalleryService{
		DB: db,
	}

	emailService := models.NewEmailService(smtpCfg)

	um := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(csrfSecure),
		csrf.Path("/"),
	)

	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	r := chi.NewRouter()

	r.Use(csrfMw)
	r.Use(um.SetUser)

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)

	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	r.Post("/signout", usersC.ProcessSignOut)

	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))

	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	usersC.Templates.UsersMe = views.Must(views.ParseFS(templates.FS, "currentuser.gohtml", "tailwind.gohtml"))
	r.Route("/users/me", func(r chi.Router) {
		r.Use(um.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	galleriesC.Templates.New = views.Must(views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.Edit = views.Must(views.ParseFS(templates.FS, "galleries/edit.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.Index = views.Must(views.ParseFS(templates.FS, "galleries/index.gohtml", "tailwind.gohtml"))
	galleriesC.Templates.Show = views.Must(views.ParseFS(templates.FS, "galleries/show.gohtml", "tailwind.gohtml"))
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)
		r.Group(func(r chi.Router) {
			r.Use(um.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images", galleriesC.UploadImage)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Printf("Starting the server on %s...\n", serverAddress)
	err = http.ListenAndServe(serverAddress, r)
	if err != nil {
		panic(err)
	}
}
