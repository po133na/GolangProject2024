package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/po133na/go-mid/pkg/shelter/validator"
)

type envelope map[string]interface{}

func (app *application) readStrings(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the URL query string.
	// If no key exists this will return an empty string "".
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value
	if s == "" {
		return defaultValue
	}

	// Otherwise, return the string
	return s
}

// readInt is a helper method on application type that reads a string value from the URL query
// string and converts it to an integer before returning. If no matching key is found then it
// returns the provided default value. If the value couldn't be converted to an integer, then we
// record an error message in the provided Validator instance, and return the default value.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the URL query string.
	s := qs.Get(key)

	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}

	// Try to convert the string value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	// Otherwise, return the converted integer value.
	return i
}

// writeJSON marshals data structure to encoded JSON response. It returns an error if there are
// any issues, else error is nil.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope,
	headers http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded JSON. Use
	// no line prefix and tab indents for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the header map
	// and add each header to the http.ResponseWriter header map. Note that it's OK if the
	// provided header map is nil. Go doesn't through an error if you try to range over (
	// or generally, read from) a nil map
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		// app.logger.PrintError(err, nil)
		return err
	}

	return nil
}
