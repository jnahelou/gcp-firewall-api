package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jnahelou/gcp-firewall-api/models"
	compute "google.golang.org/api/compute/v1"
)

// TODO move to helpers
func getVars(r *http.Request) (project, serviceProject, application, rule string) {
	vars := mux.Vars(r)
	project = vars["project"]
	serviceProject = vars["service-project"]
	application = vars["application"]
	rule = vars["rule"]
	return
}

func ListFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := getVars(r)

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

func CreateFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := getVars(r)

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
		log.Printf("[ERROR] Somes rules cannot be created : %v\n", err)
		return
	}
}

func DeleteFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := getVars(r)

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

func CreateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := getVars(r)

	var frule compute.Firewall

	err := json.NewDecoder(r.Body).Decode(&frule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("[DEBUG] Ask to create rule %s %s %s %s\n", project, serviceProject, application, rule)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.CreateFirewallRule(manager, project, serviceProject, application, rule, frule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := getVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("[DEBUG] Ask to create rule %s %s %s %s\n", project, serviceProject, application, rule)

	err = models.DeleteFirewallRule(manager, project, serviceProject, application, rule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
