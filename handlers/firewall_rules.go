package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jnahelou/gcp-firewall-api/models"
	"github.com/jnahelou/gcp-firewall-api/services"
	"github.com/sirupsen/logrus"
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

// ListFirewallRulesHandler returns a set of firewall rules
func ListFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := getVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frs, err := services.ListApplicationFirewallRules(manager, project, serviceProject, application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(frs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(res))
}

// DeleteFirewallRulesHandler delete all application firewall rules
func DeleteFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := getVars(r)
	logrus.Debugf("Ask to delete all rules for application %s in project %s/%s\n", application, project, serviceProject)

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

	err = services.DeleteApplicationFirewallRules(manager, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateFirewallRulesHandler create a set of firewall rules
func CreateFirewallRulesHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := getVars(r)

	var rules models.FirewallRuleList

	err := json.NewDecoder(r.Body).Decode(&rules)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Debugf("Ask to create rules on project '%s', for service project '%s', for application %s:\n", project, serviceProject, application)
	for _, r := range rules {
		logrus.Printf("Rule: %s\n", r.CustomName)
	}

	app := models.ApplicationRules{Project: project, ServiceProject: serviceProject, Application: application, Rules: rules}
	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.CreateApplicationFirewallRules(manager, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFirewallRuleHandler return mathing firewall rule
func GetFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := getVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	googleRule, err := services.GetFirewallRule(manager, project, serviceProject, application, rule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cast into our rule
	customRule := models.FirewallRule{
		CustomName: rule,
		Rule:       *googleRule,
	}

	res, err := json.Marshal(customRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(res))
}

// CreateFirewallRuleHandler create a given rule
func CreateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := getVars(r)
	logrus.Debugf("Ask to create rule %s %s %s %s\n", project, serviceProject, application, rule)

	// Decode given rule in order to create it
	var body compute.Firewall
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	googleRule, err := services.CreateFirewallRule(manager, project, serviceProject, application, rule, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Cast into our rule
	customRule := models.FirewallRule{
		CustomName: rule,
		Rule:       *googleRule,
	}

	res, err := json.Marshal(customRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(res))
}

// UpdateFirewallRuleHandler recreate the given firewall rule
func UpdateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := getVars(r)
	logrus.Debugf("Ask to recreate rule %s %s %s %s\n", project, serviceProject, application, rule)

	var frule compute.Firewall

	err := json.NewDecoder(r.Body).Decode(&frule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	googleRule, err := services.CreateFirewallRule(manager, project, serviceProject, application, rule, frule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cast into our rule
	customRule := models.FirewallRule{
		CustomName: rule,
		Rule:       *googleRule,
	}

	res, err := json.Marshal(customRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(res))
}

// DeleteFirewallRuleHandler delete the given firewall rule
func DeleteFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := getVars(r)
	logrus.Debugf("Ask to delete rule %s %s %s %s\n", project, serviceProject, application, rule)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.DeleteFirewallRule(manager, project, serviceProject, application, rule)
	if err != nil {
		// TODO: if rule not found, its raise an error 500. It should not. Test error type
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
