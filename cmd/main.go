package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"server/pkg/handlers"
	"syscall"
	"time"

	notifyMiddleware "server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var port = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("failed loading .env file")
	}

	path := os.Getenv("NOTIFY_DB_PATH")
	app := chi.NewRouter()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(fmt.Sprintf("failed database connection: %v", err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("failed database ping: %v", err))
	}

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app.Use(
		middleware.Logger,
		notifyMiddleware.Middleware,
	)

	filerServer := http.FileServer(http.Dir("./public"))
	app.Handle("/public/*", http.StripPrefix("/public/", filerServer))

	app.Group(func(r chi.Router) {
		app.NotFound(handlers.NewNotFoundHanlder().ServeHTTP)
		// public
		app.Get("/", handlers.NewHomeHandler(slogger).ServeHTTP)
		app.Get("/login", handlers.NewLoginHanlder(slogger).ServeHTTP)
		app.Get("/signup", handlers.NewSignupHandler(slogger).ServeHTTP)
		app.Get("/profile", handlers.NewProfileHandler(db, slogger).ServeHTTP)
		// public forms
		app.Post("/login", handlers.NewPostLoginHandler(db, slogger).ServeHTTP)
		app.Post("/signup", handlers.NewPostSignupHandler(db, slogger).ServeHTTP)
		app.Post("/logout", handlers.NewPostLogOutHandler().ServeHTTP)
		// user profile
		app.Get("/jobs", handlers.NewGetJobHandler(db, slogger).ServeHTTP)

		// admin
		app.Get("/admin", handlers.NewAdminHandler(db, slogger).ServeHTTP)
		app.Get("/admin/jobs", handlers.NewAdminJobHandler(db, slogger).ServeHTTP)
		app.Post("/admin/jobs", handlers.NewPostSearchJobHandler(db, slogger).ServeHTTP)
		app.Delete("/admin/jobs/{id}", handlers.NewDeleteAdminJobHandler(db, slogger).ServeHTTP)
		app.Get("/admin/users", handlers.NewAdminUserHandler(db, slogger).ServeHTTP)
		app.Post("/admin/users", handlers.NewPostSearchUserHandler(db, slogger).ServeHTTP)
		app.Delete("/admin/users/{id}", handlers.NewDeleteAdminUserHandler(db, slogger).ServeHTTP)
		// admin forms
		app.Get("/admin/jobs/new", handlers.NewGetAdminCreateJobHandler(db, slogger).ServeHTTP)
		app.Post("/admin/jobs/new", handlers.NewPostAdminCreateJobHandler(db, slogger).ServeHTTP)
		app.Get("/admin/jobs/edit/{id}", handlers.NewGetAdminJobEditHandler(db, slogger).ServeHTTP)
		app.Put("/admin/jobs/edit/{id}", handlers.NewPutAdminJobEditHandler(db, slogger).ServeHTTP)
		app.Get("/admin/users/new", handlers.NewGetAdminCreateUserHandler(db, slogger).ServeHTTP)
		app.Post("/admin/users/new", handlers.NewPostAdminCreateUserHandler(db, slogger).ServeHTTP)
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
		// slogger.Info("Started server on port %s", port)
		slogger.Info(fmt.Sprintf("Started sever on port %s", port))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slogger.Error(fmt.Sprintf("ListenAndServe() :%v", err))
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
		slogger.Error(fmt.Sprintf("sever shutdown failed: %v", err))
		os.Exit(1)
	}
	slogger.Info("Sever shutdown complete")
}
