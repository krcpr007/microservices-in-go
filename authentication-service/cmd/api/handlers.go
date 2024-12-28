package main

import (
	"errors"
	"net/http"
	"data"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user  
	user, err := app.Models.Users.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("Invalid Creadentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password); 
	
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid Creadentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error: false, 
		Message: "Logged In",
		Data: user,
	}

	app.writeJSON(w, http.StatusOK, payload)
}