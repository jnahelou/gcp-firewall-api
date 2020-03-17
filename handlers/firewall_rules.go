package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jnahelou/gcp-firewall-api/models"
)

// TODO move to helpers
func getVars(r *http.Request) (string, string, string) {
	vars := mux.Vars(r)
	project := vars["project"]
	serviceProject := vars["service-project"]
	application := vars["application"]

	return project, serviceProject, application
}

func ListFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application := getVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frs, err := models.ListApplicationFirewallRules(manager, project, serviceProject, application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(frs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(res))
}

func CreateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application := getVars(r)

	var rules models.FirewallRuleList

	err := json.NewDecoder(r.Body).Decode(&rules)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app := models.ApplicationRules{Project: project, ServiceProject: serviceProject, Application: application, Rules: rules}
	fmt.Printf("[DEBUG] Ask to create %v\n", app)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.CreateApplicationFirewallRules(manager, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func DeleteFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application := getVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app := models.ApplicationRules{
		Project:        project,
		ServiceProject: serviceProject,
		Application:    application,
		Rules:          nil,
	}

	err = models.DeleteApplicationFirewallRules(manager, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
