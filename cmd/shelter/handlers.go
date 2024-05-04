package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/po133na/go-mid/pkg/shelter/model"
	"github.com/po133na/go-mid/pkg/shelter/validator"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createAnimalHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID             string `json:"id"`
		Kind_Of_Animal string `json:"kind_of_animal"`
		Kind_Of_Breed  string `json:"kind_of_breed"`
		Name           string `json:"name"`
		Age            string `json:"age"`
		Description    string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	animal := &model.Animal{
		ID:             input.ID,
		Kind_Of_Animal: input.Kind_Of_Animal,
		Kind_Of_Breed:  input.Kind_Of_Breed,
		Name:           input.Name,
		Age:            input.Age,
		Description:    input.Description,
	}

	err = app.models.Animals.Insert(animal)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, animal)
}

func (app *application) getAnimalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["animalId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid animal ID")
		return
	}

	animal, err := app.models.Animals.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, animal)
}

func (app *application) getAnimalsSortedHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// param := vars["animalId"]
	var input struct {
		Kind_Of_Animal string `json:"kind_of_animal"`
		Kind_Of_Breed  string `json:"kind_of_breed"`
		// Age            int    `json:"age"`
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Kind_Of_Animal = app.readStrings(qs, "kind_of_animal", "")
	input.Kind_Of_Breed = app.readStrings(qs, "kind_of_breed", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "name", "kind_of_breed", "kind_of_animal", "age",
		// descending sort values
		"-id", "-name", "-kind_of_breed", "-kind_of_animal", "-age",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	animals, metadata, err := app.models.Animals.GetSort(input.Kind_Of_Breed, input.Kind_Of_Animal, input.Filters)
	if err != nil {
		fmt.Println("We are in search animals handler", "\nanimal: ", input.Kind_Of_Animal, "\nbreed:", input.Kind_Of_Breed, "\n", input.Filters)
		fmt.Print("\nError: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"animals": animals, "metadata": metadata}, nil)
}

func (app *application) updateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["animalId"] // CHECK HERE FOR ERRORS

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid animal ID")
		return
	}

	animal, err := app.models.Animals.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Name *string `json:"name"`
		Age  *string `json:"age"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		animal.Name = *input.Name
	}

	if input.Age != nil {
		animal.Age = *input.Age
	}
	err = app.models.Animals.Update(animal)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) deleteAnimalHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["animalId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Animal ID")
		return
	}

	err = app.models.Animals.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
