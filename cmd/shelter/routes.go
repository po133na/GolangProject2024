package main

import (
	"net/http"

	//new
	"github.com/gorilla/mux"
	//new
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	// Animal Singleton
	v1.HandleFunc("/animals", app.createAnimalHandler).Methods("POST")
	v1.HandleFunc("/animals/{animalId:[0-9]+}", app.getAnimalHandler).Methods("GET")
	v1.HandleFunc("/animals/sort", app.getAnimalsSortedHandler).Methods("GET")
	v1.HandleFunc("/animals/{animalId:[0-9]+}", app.updateAnimalHandler).Methods("PUT")
	v1.HandleFunc("/animals/{animalId:[0-9]+}", app.requirePermissions("animals:read", app.deleteAnimalHandler)).Methods("DELETE")

	users1 := r.PathPrefix("/api/v1").Subrouter()

	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	return app.authenticate(r)
}
