package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type FirewallRule struct {
	Rule       compute.Firewall `json:firewall`
	CustomName string           `json:customName`
}
type FirewallRuleList []FirewallRule

type FirewallRuleManager interface {
	ListFirewallRules(gcp_project string) ([]*compute.Firewall, error)
	GetFirewallRule(project, name string) (*compute.Firewall, error)
	CreateFirewallRule(project string, rule *compute.Firewall) error
	UpdateFirewallRule(project string, rule *compute.Firewall) error
	DeleteFirewallRule(project, name string) error
}

/*
 * FirewallRuleClient provides primitives to collect rules from Google Cloud Platform
 * FirewallRuleClient implements FirewallRuleManager
 *
 */
type FirewallRuleClient struct {
	computeService *compute.Service
}

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

func (f *FirewallRuleClient) ListFirewallRules(project string) ([]*compute.Firewall, error) {
	ctx := context.Background()

	req := f.computeService.Firewalls.List(project)

	var firewall_rule_list []*compute.Firewall

	if err := req.Pages(ctx, func(page *compute.FirewallList) error {
		for _, firewall := range page.Items {
			firewall_rule_list = append(firewall_rule_list, firewall)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return firewall_rule_list, nil
}

func (f *FirewallRuleClient) GetFirewallRule(project, name string) (*compute.Firewall, error) {
	ctx := context.Background()
	return f.computeService.Firewalls.Get(project, name).Context(ctx).Do()

}

func (f *FirewallRuleClient) CreateFirewallRule(project string, rule *compute.Firewall) error {
	ctx := context.Background()

	resp, err := f.computeService.Firewalls.Insert(project, rule).Context(ctx).Do()
	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] CreateFirewallRule result : %v\n", resp)

	_, err = f.GetFirewallRule(project, rule.Name)
	return err
}
func (f *FirewallRuleClient) UpdateFirewallRule(project string, rule *compute.Firewall) error {
	_, err := f.GetFirewallRule(project, rule.Name)
	return err
}

func (f *FirewallRuleClient) DeleteFirewallRule(project string, name string) error {
	ctx := context.Background()

	resp, err := f.computeService.Firewalls.Delete(project, name).Context(ctx).Do()
	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] DeleteFirewallRule result : %v\n", resp)

	return err
}

/*
 *
 * ApplicationRules defines a set of rules apply for an Application in a GCP Project
 *
 */
type ApplicationRules struct {
	Project        string
	ServiceProject string
	Application    string
	Rules          FirewallRuleList
}

func ListApplicationFirewallRules(manager FirewallRuleManager, project, serviceProject, application string) (*ApplicationRules, error) {
	fmt.Printf("[DEBUG] Manager will list rules for project %s\n", project)

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

func CreateApplicationFirewallRules(manager FirewallRuleManager, appRules ApplicationRules) error {
	project := appRules.Project
	serviceProject := appRules.ServiceProject
	application := appRules.Application

	var err error
	for _, rule := range appRules.Rules {
		rule.Rule.Name = fmt.Sprintf("%s-%s-%s", serviceProject, application, rule.CustomName)
		fmt.Printf("[DEBUG] Manager will create %s on %s\n", rule.Rule.Name, project)
		err = manager.CreateFirewallRule(project, &rule.Rule)
	}

	return err
}

func DeleteApplicationFirewallRules(manager FirewallRuleManager, appRules ApplicationRules) error {
	project := appRules.Project
	serviceProject := appRules.ServiceProject
	application := appRules.Application

	rules, err := ListApplicationFirewallRules(manager, project, serviceProject, application)
	if err != nil {
		log.Printf("[ERROR] Unable to list rules during DeleteApplicationFirewallRules. Got error : %v", err)
		return err
	}
	for _, rule := range rules.Rules {
		name := fmt.Sprintf("%s-%s", serviceProject, application)
		fmt.Printf("[DEBUG] Manager will delete %s on %s. Content %v\n", name, project, rule)

		err = manager.DeleteFirewallRule(project, rule.Rule.Name)
		if err != nil {
			fmt.Printf("[ERROR] Unable to delete rule %s. Got error : %v\n", rule.Rule.Name, err)
			return err
		}
	}

	return err
}
