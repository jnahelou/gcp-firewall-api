package services

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/adeo/iwc-gcp-firewall-api/models"
	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
)

func init() {
	logrus.SetOutput(ioutil.Discard)
}

// FirewallRuleDummyClient provides primitives to collect rules from in-memory rules list
type FirewallRuleDummyClient struct {
	Rules map[string][]*compute.Firewall
}

func NewFirewallRuleDummyClient() (*FirewallRuleDummyClient, error) {
	manager := FirewallRuleDummyClient{}
	manager.Rules = make(map[string][]*compute.Firewall)
	return &manager, nil
}

func (f *FirewallRuleDummyClient) ListFirewallRule(project string) ([]*compute.Firewall, error) {
	if value, ok := f.Rules[project]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("Project not found")
}

func (f *FirewallRuleDummyClient) GetFirewallRule(project, name string) (*compute.Firewall, error) {
	for _, rule := range f.Rules[project] {
		if rule.Name == name {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("Rule not found")
}

func (f *FirewallRuleDummyClient) CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error) {
	for _, r := range f.Rules[project] {
		if r.Name == rule.Name {
			return nil, fmt.Errorf("Rule already exists")
		}
	}

	f.Rules[project] = append(f.Rules[project], rule)
	return rule, nil
}

func (f *FirewallRuleDummyClient) DeleteFirewallRule(project, name string) error {
	rules := f.Rules[project]
	for i, rule := range rules {
		if rule.Name == name {
			// Delete matching route
			rules[i] = rules[len(rules)-1]
			f.Rules[project] = rules[:len(rules)-1]
			return nil
		}
	}
	return fmt.Errorf("Rule not found")
}

func TestCreateFirewallRule(t *testing.T) {
	manager, _ := NewFirewallRuleDummyClient()
	project := "dummy-project"
	serviceProject := "dummy-service_project"
	application := "dummy-application"
	var rules models.FirewallRules
	rule := models.FirewallRule{
		Rule:       compute.Firewall{Name: "remote", Network: "global/networks/default", Allowed: []*compute.FirewallAllowed{&compute.FirewallAllowed{Ports: []string{"22", "3389"}, IPProtocol: "TCP"}}},
		CustomName: "allow-tcp-22-3389",
	}
	rules = append(rules, rule)

	// Create dummy rule
	for _, rule := range rules {
		_, err := CreateFirewallRule(manager, project, serviceProject, application, rule.CustomName, rule.Rule)
		if err != nil {
			t.Fatalf("Something wrong during rule creation. Got error %v\n", err)
		}
	}

	// Test if expected
	expected := fmt.Sprintf("%s-%s-%s", serviceProject, application, rule.CustomName)
	if manager.Rules[project][0].Name != expected {
		t.Errorf("Name don't match format. Got %s, expected %s\n", manager.Rules[project][0].Name, expected)

	}

	// Inster existing rule should trigger error
	_, err := CreateFirewallRule(manager, project, serviceProject, application, rule.CustomName, rule.Rule)
	if err == nil {
		t.Errorf("Expected error during insert if rule already exists")
	}

}

func TestListFirewallRule(t *testing.T) {
	// Add dummy content
	manager, _ := NewFirewallRuleDummyClient()

	project := "kubernetes-host-project"
	serviceProjects := []string{"kubernetes-demo", "kubernetes-training"}
	applications := []string{"the-hard-way", "the-easy-way"}

	// Create 4 dummy rules
	for _, serviceProject := range serviceProjects {
		for _, application := range applications {
			name := fmt.Sprintf("%s-%s-%s", serviceProject, application, "allow-external")
			rule := compute.Firewall{Name: name, Network: "global/networks/default", Allowed: []*compute.FirewallAllowed{&compute.FirewallAllowed{Ports: []string{"22", "6443"}, IPProtocol: "TCP"}}}
			manager.Rules[project] = append(manager.Rules[project], &rule)
		}
	}

	// Ask for non-existing project
	_, err := ListFirewallRule(manager, "non-existing-project", serviceProjects[0], applications[0])
	if err == nil {
		t.Fatalf("Expected error during ListApplicationFirewallRules on a non-existing project")
	}

	// Ask for one application in a random project
	applicationRule, err := ListFirewallRule(manager, project, serviceProjects[0], applications[0])
	if err != nil {
		t.Fatalf("Something wrong during ListApplicationFirewallRules. Got error : %v\n", err)
	}

	if len(applicationRule.Rules) != 1 {
		t.Errorf("Wrong rules count for project %s and application %s. Got %d expected %d", serviceProjects[0], applications[0], len(applicationRule.Rules), 1)
	}
}

func TestDeleteApplicationFirewallRules(t *testing.T) {
	// Add dummy content
	manager, _ := NewFirewallRuleDummyClient()
	project := "nginx-host-project"
	serviceProject := "nginx-demo"
	application := "front"
	ruleCustomName := "allow-publicly"
	name := fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleCustomName)
	gRule := compute.Firewall{Name: name, Network: "global/networks/default", Allowed: []*compute.FirewallAllowed{&compute.FirewallAllowed{Ports: []string{"80", "443"}, IPProtocol: "TCP"}}}
	manager.Rules[project] = append(manager.Rules[project], &gRule)

	// Ask to delete a rule
	err := DeleteFirewallRule(manager, project, serviceProject, application, ruleCustomName)
	if err != nil {
		t.Fatalf("Unexpected error during Delete. Got %v\n", err)
	}

	// Verify empty rules
	apprules, err := ListFirewallRule(manager, project, serviceProject, application)
	if err != nil {
		t.Fatalf("Unexpected error during Delete. Got %v\n", err)
	}
	expected := 0
	if len(apprules.Rules) != expected {
		t.Errorf("Bad rules count. Got %d expected %d", len(apprules.Rules), expected)
	}

	// Try to delete on non-existing project
	err = DeleteFirewallRule(manager, project, serviceProject, application, ruleCustomName)
	if err == nil {
		t.Fatalf("Expected error during Delete on non existing project. Got %v\n", err)
	}
}
