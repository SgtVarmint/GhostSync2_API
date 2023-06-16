package authentication

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/SgtVarmint/GhostSync2/config"
)

type ValidAuth struct {
	Status		string
}

func Authenticate(w http.ResponseWriter, r *http.Request, config config.Config) {
	queryStrings := r.URL.Query()

	var payload *ValidAuth
	if (queryStrings.Get("accessCode") == config.AccessCode) {
		payload = &ValidAuth { Status: config.AuthenticatedCode }
	} else {
		payload = &ValidAuth { Status: "denied" }
	}

	json, err := json.Marshal(payload)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}