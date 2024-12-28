package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)



type RequestPayload struct {
	Action string `json:"action"`
	Auth string `json:"auth,omitempty`
}

type AuthPayload struct {
	Email string `json:"email`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}


func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJSON(w,r,&requestPayload); 
	
	if err != nil {
		app.errorJSON(w,err)
	}

	switch requestPayload.Action {
		case "auth":
			var authPayload AuthPayload
			err := json.Unmarshal([]byte(requestPayload.Auth), &authPayload)
			if err != nil {
				app.errorJSON(w, errors.New("Invalid auth payload"))
				return
			}
			app.authenticate(w, authPayload)

		default:
			app.errorJSON(w,errors.New("Invalid action"))
	}
}

func(app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _  := json.MarshalIndent(a,"","\t")


	request, err := http.NewRequest("POST","http://authentication-service/authenticate",bytes.NewBuffer(jsonData))

	client := &http.Client{}

	response, err:= client.Do(request)

	if err != nil {
		app.errorJSON(w,err)
		return; 
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w,errors.New("Invalid Credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w,errors.New("Error calling auth service"))
		return
	}


	// create a variable we will send response.body into 

	var jsonFormService jsonResponse 

	err = json.NewDecoder(response.Body).Decode(&jsonFormService)

	if err != nil {
		app.errorJSON(w,err)
		return; 
	}

	if jsonFormService.Error { 
		app.errorJSON(w,errors.New(jsonFormService.Message))
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFormService.Data 

	app.writeJSON(w,http.StatusOK,payload)
}