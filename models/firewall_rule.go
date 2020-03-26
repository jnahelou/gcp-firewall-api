package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

// FirewallRule descibe a firewall rule
type FirewallRule struct {
	Rule       compute.Firewall `json:"rule"`
	CustomName string           `json:"custom_name"`
}

// FirewallRuleList describe a set of firewall rull
type FirewallRuleList []FirewallRule

// FirewallRuleManager contains methods to manage firewall rules
type FirewallRuleManager interface {
	ListFirewallRules(project string) ([]*compute.Firewall, error)
	GetFirewallRule(project, name string) (*compute.Firewall, error)
	CreateFirewallRule(project string, rule *compute.Firewall) error
	UpdateFirewallRule(project string, rule *compute.Firewall) error
	DeleteFirewallRule(project, name string) error
}

// FirewallRuleClient provides primitives to collect rules from Google Cloud Platform. Implements FirewallRuleManager
type FirewallRuleClient struct {
	computeService *compute.Service
}

// NewFirewallRuleClient FirewallRuleClient contructor
func NewFirewallRuleClient() (*FirewallRuleClient, error) {
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	computeService, err := compute.New(c)
	if err != nil {
		return nil, err
	}

	manager := FirewallRuleClient{}
	manager.computeService = computeService
	return &manager, err
}

// ListFirewallRules returns given project's firewall rule
func (f *FirewallRuleClient) ListFirewallRules(project string) ([]*compute.Firewall, error) {
	ctx := context.Background()

	req := f.computeService.Firewalls.List(project)

	var firewallRuleList []*compute.Firewall

	if err := req.Pages(ctx, func(page *compute.FirewallList) error {
		for _, firewall := range page.Items {
			firewallRuleList = append(firewallRuleList, firewall)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return firewallRuleList, nil
}

// GetFirewallRule returns firewall rule matching given project and name
func (f *FirewallRuleClient) GetFirewallRule(project, name string) (*compute.Firewall, error) {
	ctx := context.Background()
	return f.computeService.Firewalls.Get(project, name).Context(ctx).Do()

}

// CreateFirewallRule create given firewall rule on given project
func (f *FirewallRuleClient) CreateFirewallRule(project string, rule *compute.Firewall) error {
	ctx := context.Background()

	resp, err := f.computeService.Firewalls.Insert(project, rule).Context(ctx).Do()
	if err != nil {
		return err
	}

	logrus.Debugf("CreateFirewallRule result : %v\n", resp)

	_, err = f.GetFirewallRule(project, rule.Name)
	return err
}

// UpdateFirewallRule update given firewall rule in given project
func (f *FirewallRuleClient) UpdateFirewallRule(project string, rule *compute.Firewall) error {
	_, err := f.GetFirewallRule(project, rule.Name)
	return err
}

// DeleteFirewallRule delete firewall rule matching given project and name
func (f *FirewallRuleClient) DeleteFirewallRule(project string, name string) error {
	ctx := context.Background()

	resp, err := f.computeService.Firewalls.Delete(project, name).Context(ctx).Do()
	if err != nil {
		return err
	}

	logrus.Debugf("DeleteFirewallRule result : %v\n", resp)

	return err
}

// ApplicationRules defines a set of rules apply for an Application in a GCP Project
type ApplicationRules struct {
	Project        string
	ServiceProject string
	Application    string
	Rules          FirewallRuleList
}

// ListApplicationFirewallRules returns a set of firewall rules related to an application
func ListApplicationFirewallRules(manager FirewallRuleManager, project, serviceProject, application string) (*ApplicationRules, error) {
	logrus.Debugf("Manager will list rules for project %s\n", project)

	rules, err := manager.ListFirewallRules(project)
	if err != nil {
		return nil, err
	}

	var applicationRules FirewallRuleList
	for _, rule := range rules {
		prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
		if strings.HasPrefix(rule.Name, prefix) {
			customName := rule.Name[len(prefix):]
			applicationRules = append(applicationRules, FirewallRule{Rule: *rule, CustomName: customName})
		}
	}

	return &ApplicationRules{Project: project, ServiceProject: serviceProject, Application: application, Rules: applicationRules}, nil
}

// CreateFirewallRule create given firewall rule on given project
func CreateFirewallRule(manager FirewallRuleManager, project string, serviceProject string, application string, ruleName string, rule compute.Firewall) error {
	rule.Name = fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)
	logrus.Debugf("Manager will create %s on %s\n", rule.Name, project)
	return manager.CreateFirewallRule(project, &rule)
}

// CreateApplicationFirewallRules create all provided firewall rules
func CreateApplicationFirewallRules(manager FirewallRuleManager, appRules ApplicationRules) error {
	project := appRules.Project
	serviceProject := appRules.ServiceProject
	application := appRules.Application

	var errs []error
	for _, rule := range appRules.Rules {
		ruleName := rule.CustomName
		err := CreateFirewallRule(manager, project, serviceProject, application, ruleName, rule.Rule)
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
func DeleteFirewallRule(manager FirewallRuleManager, project, serviceProject, application, customName string) error {
	ruleName := fmt.Sprintf("%s-%s-%s", serviceProject, application, customName)
	logrus.Debugf("Manager will delete %s on %s.\n", ruleName, project)
	return manager.DeleteFirewallRule(project, ruleName)
}

// DeleteApplicationFirewallRules delete given set of application's firewall rule
func DeleteApplicationFirewallRules(manager FirewallRuleManager, appRules ApplicationRules) error {
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
