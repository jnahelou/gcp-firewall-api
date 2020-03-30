package handlers

import (
	"fmt"
	"net/http"
)

// HealthCheckHandler ensure application is runniing properly
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"ping": "pong"}`)
}
