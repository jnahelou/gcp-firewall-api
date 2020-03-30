package models

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

// FirewallRule descibe a firewall rule
type FirewallRule struct {
	Rule       compute.Firewall `json:"item"`
	CustomName string           `json:"custom_name"`
}

// FirewallRules describe a set of firewall rule
type FirewallRules []FirewallRule

// ApplicationRule describe and end-user response
type ApplicationRule struct {
	Project        string        `json:"project"`
	ServiceProject string        `json:"service_project"`
	Application    string        `json:"application"`
	Rules          FirewallRules `json:"data"`
}

// FirewallRuleManager contains methods to manage firewall rules
type FirewallRuleManager interface {
	ListFirewallRule(project string) ([]*compute.Firewall, error)
	GetFirewallRule(project, name string) (*compute.Firewall, error)
	CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error)
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

// ListFirewallRule returns given project's firewall rule
func (f *FirewallRuleClient) ListFirewallRule(project string) ([]*compute.Firewall, error) {
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
func (f *FirewallRuleClient) CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error) {
	_, err := f.computeService.Firewalls.Insert(project, rule).Context(context.Background()).Do()
	if err != nil {
		return nil, err
	}

	return f.GetFirewallRule(project, rule.Name)
}

// DeleteFirewallRule delete firewall rule matching given project and name
func (f *FirewallRuleClient) DeleteFirewallRule(project string, name string) error {
	_, err := f.computeService.Firewalls.Delete(project, name).Context(context.Background()).Do()
	return err
}
