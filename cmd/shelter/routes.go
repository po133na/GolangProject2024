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
	//User Singleton
	users1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/users/login", app.createAuthenticationTokenHandler).Methods("POST")

	//Shelter Singleton
	sh1 := r.PathPrefix("/api/v1").Subrouter()

	sh1.HandleFunc("/shelters", app.createShelterHandler).Methods("POST")
	sh1.HandleFunc("/shelters/{shelterId:[0-9]+}", app.getShelterHandler).Methods("GET")
	sh1.HandleFunc("/shelters/sort", app.getSheltersSortedHandler).Methods("GET")
	sh1.HandleFunc("/shelters/{shelterId:[0-9]+}", app.updateShelterHandler).Methods("PUT")
	sh1.HandleFunc("/shelters/{shelterId:[0-9]+}", app.requirePermissions("shelters:read", app.deleteShelterHandler)).Methods("DELETE")
	//get a list of employees from shelter
	sh1.HandleFunc("/shelters/{shelterId}/employees", app.getEmployeesSortedHandler).Methods("GET")
	//Employee Singleton
	em1 := r.PathPrefix("/api/v1").Subrouter()

	em1.HandleFunc("/employee", app.createEmployeeHandler).Methods("POST")
	em1.HandleFunc("/employees/{employeeId:[0-9]+}", app.getEmployeeHandler).Methods("GET")
	em1.HandleFunc("/employees/sort", app.getEmployeesSortedHandler).Methods("GET")
	em1.HandleFunc("/employees/{employeeId:[0-9]+}", app.updateEmployeeHandler).Methods("PUT")
	em1.HandleFunc("/employees/{employeeId:[0-9]+}", app.requirePermissions("employees:read", app.deleteEmployeeHandler)).Methods("DELETE")
	return app.authenticate(r)
}
