package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jnahelou/gcp-firewall-api/helpers"
	"github.com/jnahelou/gcp-firewall-api/models"
	"github.com/jnahelou/gcp-firewall-api/services"
	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// ListFirewallRuleHandler returns a set of firewall rules
func ListFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, _ := helpers.GetMuxVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	applicationRule, err := services.ListFirewallRule(manager, project, serviceProject, application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(applicationRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(res))
}

// GetFirewallRuleHandler return mathing firewall rule
func GetFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := helpers.GetMuxVars(r)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	applicationRule, err := services.GetFirewallRule(manager, project, serviceProject, application, rule)

	// Handle a proper 404
	if value, ok := err.(*googleapi.Error); ok {
		w.WriteHeader(value.Code)
		fmt.Fprint(w, models.NewGoogleApplicationError(value).JSON())
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(applicationRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(res))
}

// CreateFirewallRuleHandler create a given rule
func CreateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := helpers.GetMuxVars(r)
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

	applicationRule, err := services.CreateFirewallRule(manager, project, serviceProject, application, rule, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(applicationRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(res))
}

// DeleteFirewallRuleHandler delete the given firewall rule
func DeleteFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	project, serviceProject, application, rule := helpers.GetMuxVars(r)
	logrus.Debugf("Ask to delete rule %s %s %s %s\n", project, serviceProject, application, rule)

	manager, err := models.NewFirewallRuleClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.DeleteFirewallRule(manager, project, serviceProject, application, rule)

	// Handle a proper 404
	if value, ok := err.(*googleapi.Error); ok {
		w.WriteHeader(value.Code)
		fmt.Fprint(w, models.NewGoogleApplicationError(value).JSON())
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
