package config

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// Struct to hold main variables for this application.
// Routes all have access to an instance of this struct.
type App struct {
	Negroni *negroni.Negroni
	Router  *mux.Router
	Render  *render.Render
}

// This function is called from main.go and from the tests
// to create a new application.
func NewApp(root string) *App {

	CheckEnv()

	// Use negroni for middleware
	ne := negroni.New()

	// Use gorilla/mux for routing
	ro := mux.NewRouter()

	// Use Render for template. Pass in path to templates folder
	// as well as asset helper functions.
	re := render.New(render.Options{
		Directory:  filepath.Join(root, "templates"),
		Layout:     "layouts/layout",
		Extensions: []string{".html"},
		Funcs:      []template.FuncMap{AssetHelpers(root)},
	})

	// Establish connection to DB as specified in database.go
	db := NewDB()
	models.New(db)

	// Add middleware to the stack
	ne.Use(negroni.NewRecovery())
	ne.Use(negroni.NewLogger())
	ne.Use(NewAssetHeaders())
	ne.Use(negroni.NewStatic(http.Dir("public")))
	ne.UseHandler(ro)

	// Return a new App struct with all these things.
	return &App{ne, ro, re}
}
