package services

import (
	"fmt"
	"strings"

	"github.com/adeo/iwc-gcp-firewall-api/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

// ListFirewallRule returns a set of firewall rules related to an application
func ListFirewallRule(manager models.FirewallRuleManager, project, serviceProject, application string) (*models.ApplicationRule, error) {
	logrus.Debugf("Manager will list rules for project %s\n", project)

	// List all firewall rule in given project
	gRules, err := manager.ListFirewallRule(project)
	if err != nil {
		return nil, err
	}

	// Create reponse
	endUserResult := models.ApplicationRule{
		Application:    application,
		Project:        project,
		ServiceProject: serviceProject,
	}

	var endUserResultRules models.FirewallRules

	// For each obtains Google rules
	prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
	for _, gRule := range gRules {
		// Filter with managed rules with this application
		if strings.HasPrefix(gRule.Name, prefix) {
			customName := gRule.Name[len(prefix):]
			endUserResultRules = append(endUserResultRules, models.FirewallRule{
				Rule:       *gRule,
				CustomName: customName,
			})
		}
	}

	endUserResult.Rules = endUserResultRules
	return &endUserResult, nil
}

// CreateFirewallRule create given firewall rule on given project
func CreateFirewallRule(manager models.FirewallRuleManager, project string, serviceProject string, application string, ruleName string, rule compute.Firewall) (*models.ApplicationRule, error) {
	rule.Name = fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)
	logrus.Debugf("Manager will create %s on %s\n", rule.Name, project)
	gRule, err := manager.CreateFirewallRule(project, &rule)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
	createdRule := models.FirewallRule{
		Rule:       *gRule,
		CustomName: gRule.Name[len(prefix):],
	}

	return &models.ApplicationRule{
		Application:    application,
		Project:        project,
		ServiceProject: serviceProject,
		Rules:          models.FirewallRules{createdRule},
	}, nil
}

// GetFirewallRule return matching firewall rule
func GetFirewallRule(manager models.FirewallRuleManager, project string, serviceProject string, application string, ruleName string) (*models.ApplicationRule, error) {
	logrus.Debugf("Searching rule mathing project '%s', service project '%s', application '%s' and name '%s'", project, serviceProject, application, ruleName)
	n := fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)
	gRule, err := manager.GetFirewallRule(project, n)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
	createdRule := models.FirewallRule{
		Rule:       *gRule,
		CustomName: gRule.Name[len(prefix):],
	}

	return &models.ApplicationRule{
		Application:    application,
		Project:        project,
		ServiceProject: serviceProject,
		Rules:          models.FirewallRules{createdRule},
	}, nil
}

// DeleteFirewallRule delete firewall rule mathing project, service project, application name and rule name
func DeleteFirewallRule(manager models.FirewallRuleManager, project, serviceProject, application, customName string) error {
	ruleName := fmt.Sprintf("%s-%s-%s", serviceProject, application, customName)
	logrus.Debugf("Manager will delete %s on %s.\n", ruleName, project)
	return manager.DeleteFirewallRule(project, ruleName)
}
