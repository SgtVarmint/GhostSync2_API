package authentication

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/SgtVarmint/GhostSync2/config"
)

type ValidAuth struct {
	Status		string
}

func Authenticate(r *http.Request) ([]byte){
	config, err := config.ParseConfig()
	if err != nil {
		log.Println(err)
	}

	queryStrings := r.URL.Query()

	var payload *ValidAuth
	if (queryStrings.Get("accessCode") == config.AccessCode) {
		payload = &ValidAuth { Status: config.AuthenticatedCode }
	} else {
		payload = &ValidAuth { Status: "denied" }
	}

	json, err := json.Marshal(payload)
	if err != nil {
		log.Print(err)
	}

	return json
}

func SendAuthResponse(w http.ResponseWriter, authJson []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(authJson)
}