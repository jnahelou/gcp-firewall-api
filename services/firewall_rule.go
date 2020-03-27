package services

import (
	"fmt"
	"strings"

	"github.com/jnahelou/gcp-firewall-api/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

// ListApplicationFirewallRules returns a set of firewall rules related to an application
func ListApplicationFirewallRules(manager models.FirewallRuleManager, project, serviceProject, application string) (*models.ApplicationRules, error) {
	logrus.Debugf("Manager will list rules for project %s\n", project)

	rules, err := manager.ListFirewallRules(project)
	if err != nil {
		return nil, err
	}

	var applicationRules models.FirewallRuleList
	for _, rule := range rules {
		prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
		if strings.HasPrefix(rule.Name, prefix) {
			customName := rule.Name[len(prefix):]
			applicationRules = append(applicationRules, models.FirewallRule{Rule: *rule, CustomName: customName})
		}
	}

	return &models.ApplicationRules{Project: project, ServiceProject: serviceProject, Application: application, Rules: applicationRules}, nil
}

// CreateFirewallRule create given firewall rule on given project
func CreateFirewallRule(manager models.FirewallRuleManager, project string, serviceProject string, application string, ruleName string, rule compute.Firewall) (*compute.Firewall, error) {
	rule.Name = fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)
	logrus.Debugf("Manager will create %s on %s\n", rule.Name, project)
	return manager.CreateFirewallRule(project, &rule)
}

// GetFirewallRule return matching firewall rule
func GetFirewallRule(manager models.FirewallRuleManager, project string, serviceProject string, application string, ruleName string) (*compute.Firewall, error) {
	logrus.Debugf("Searching rule mathing project '%s', service project '%s', application '%s' and name '%s'", project, serviceProject, application, ruleName)
	n := fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)
	return manager.GetFirewallRule(project, n)
}

// CreateApplicationFirewallRules create all provided firewall rules
func CreateApplicationFirewallRules(manager models.FirewallRuleManager, appRules models.ApplicationRules) error {
	project := appRules.Project
	serviceProject := appRules.ServiceProject
	application := appRules.Application

	var errs []error
	for _, rule := range appRules.Rules {
		ruleName := rule.CustomName
		_, err := CreateFirewallRule(manager, project, serviceProject, application, ruleName, rule.Rule)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"go-error": err,
			}).Error("Error creating firewall rules")
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("At least one rule cannot be created. Error : %v", errs)
	}
	return nil
}

// DeleteFirewallRule delete firewall rule mathing project, service project, application name and rule name
func DeleteFirewallRule(manager models.FirewallRuleManager, project, serviceProject, application, customName string) error {
	ruleName := fmt.Sprintf("%s-%s-%s", serviceProject, application, customName)
	logrus.Debugf("Manager will delete %s on %s.\n", ruleName, project)
	return manager.DeleteFirewallRule(project, ruleName)
}

// DeleteApplicationFirewallRules delete given set of application's firewall rule
func DeleteApplicationFirewallRules(manager models.FirewallRuleManager, appRules models.ApplicationRules) error {
	project := appRules.Project
	serviceProject := appRules.ServiceProject
	application := appRules.Application

	rules, err := ListApplicationFirewallRules(manager, project, serviceProject, application)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"go-error": err,
		}).Error("Error listing firewall rules")
		return err
	}

	var errs []error
	for _, rule := range rules.Rules {
		err := DeleteFirewallRule(manager, project, serviceProject, application, rule.CustomName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"go-error": err,
			}).Error("Error deletinh firewall rules")
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("At least one rule cannot be destroyed. Error : %v", errs)
	}

	return nil
}
