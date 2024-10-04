package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/pkg/handlers"
	"syscall"
	"time"

	mware "server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var port = ":8080"

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("failed loading .env file")
	// }

	path := os.Getenv("NOTIFY_DB_PATH")
	app := chi.NewRouter()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed database connection %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed database ping! %v", err)
	}

	app.Use(
		middleware.Logger,
		mware.Middleware,
	)

	filerServer := http.FileServer(http.Dir("./public"))
	app.Handle("/public/*", http.StripPrefix("/public/", filerServer))

	app.Group(func(r chi.Router) {
		app.NotFound(handlers.NewNotFoundHanlder().ServeHTTP)
		// public
		app.Get("/", handlers.NewHomeHandler().ServeHTTP)
		app.Get("/login", handlers.NewLoginHanlder().ServeHTTP)
		app.Get("/signup", handlers.NewSignupHandler().ServeHTTP)
		app.Get("/profile", handlers.NewProfileHandler(db).ServeHTTP)
		// public forms
		app.Post("/login", handlers.NewPostLoginHandler(db).ServeHTTP)
		app.Post("/signup", handlers.NewPostSignupHandler(db).ServeHTTP)
		app.Post("/logout", handlers.NewPostLogOutHandler().ServeHTTP)
		// user profile
		app.Get("/jobs", handlers.NewGetJobHandler(db).ServeHTTP)

		// admin
		app.Get("/admin", handlers.NewAdminHandler(db).ServeHTTP)
		app.Get("/admin/jobs", handlers.NewAdminJobHandler(db).ServeHTTP)
		app.Post("/admin/jobs", handlers.NewPostSearchJobHandler(db).ServeHTTP)
		app.Delete("/admin/jobs/{id}", handlers.NewDeleteAdminJobHandler(db).ServeHTTP)
		app.Get("/admin/users", handlers.NewAdminUserHandler(db).ServeHTTP)
		app.Post("/admin/users", handlers.NewPostSearchUserHandler(db).ServeHTTP)
		app.Delete("/admin/users/{id}", handlers.NewDeleteAdminUserHandler(db).ServeHTTP)
		// admin forms
		app.Get("/admin/jobs/new", handlers.NewGetAdminCreateJobHandler(db).ServeHTTP)
		app.Post("/admin/jobs/new", handlers.NewPostAdminCreateJobHandler(db).ServeHTTP)
		app.Get("/admin/jobs/edit/{id}", handlers.NewGetAdminJobEditHandler(db).ServeHTTP)
		app.Put("/admin/jobs/edit/{id}", handlers.NewPutAdminJobEditHandler(db).ServeHTTP)
		app.Get("/admin/users/new", handlers.NewGetAdminCreateUserHandler(db).ServeHTTP)
		app.Post("/admin/users/new", handlers.NewPostAdminCreateUserHandler(db).ServeHTTP)
		// app.Get("/admin/users/edit/{id}", handlers.NewGetAdminUserEditHandler(db).ServeHTTP)
		// app.Put("/admin/users/edit/{id}", handlers.NewPutAdminUserEditHandler(db).ServeHTTP)
	})

	// http.ListenAndServe(port, app)

	srv := &http.Server{
		Addr:    port,
		Handler: app,
	}

	// creat channel to listen for os signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// start server in groutine
	go func() {
		log.Printf("Started server on port %s", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
			os.Exit(1)
		}
	}()

	// wait for a signal shutdown
	<-stop
	log.Println("Shutdown signal received")

	// create context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed %v", err.Error())
		os.Exit(1)
	}
	log.Println("Sever shutdown complete")
}
