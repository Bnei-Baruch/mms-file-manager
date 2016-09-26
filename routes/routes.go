package routes

import (
	"net/http"

	"github.com/Bnei-Baruch/mms-file-manager/config"
)

var App *config.App

func Setup(app *config.App) {

	App = app

	// It's important that this is before your catch-all route ("/")
	// api := App.Router.PathPrefix("/api/v1/").Subrouter()
	// api.HandleFunc("/users", GetUsersHandler).Methods("GET")
	// Optional: Use a custom 404 handler for our API paths.
	// api.NotFoundHandler = JSONNotFound

	// Serve static assets directly.
	App.Router.PathPrefix("/assets/dist/").Handler(
		http.StripPrefix("/assets/dist/", http.FileServer(http.Dir("assets/dist"))))

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	App.Router.PathPrefix("/").HandlerFunc(IndexHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/dist/index.html")
}
