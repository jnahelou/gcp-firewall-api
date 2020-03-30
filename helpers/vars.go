package helpers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// GetMuxVars return query vars from given request
func GetMuxVars(r *http.Request) (project, serviceProject, application, rule string) {
	vars := mux.Vars(r)
	project = vars["project"]
	serviceProject = vars["service_project"]
	application = vars["application"]
	rule = vars["rule"]
	return
}
