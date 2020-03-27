package models

import (
	"context"

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
	CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error)
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
func (f *FirewallRuleClient) CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error) {
	ctx := context.Background()

	resp, err := f.computeService.Firewalls.Insert(project, rule).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	logrus.Debugf("CreateFirewallRule result : %v\n", resp)

	return f.GetFirewallRule(project, rule.Name)
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
