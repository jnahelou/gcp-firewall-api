package models

import (
	"fmt"
	"testing"

	compute "google.golang.org/api/compute/v1"
)

/*
 * FirewallRuleDummyClient provides primitives to collect rules from in-memory rules list
 *
 */
type FirewallRuleDummyClient struct {
	Rules map[string][]*compute.Firewall
}

func NewFirewallRuleDummyClient() (*FirewallRuleDummyClient, error) {
	manager := FirewallRuleDummyClient{}
	manager.Rules = make(map[string][]*compute.Firewall)
	return &manager, nil
}

func (f *FirewallRuleDummyClient) ListFirewallRules(project string) ([]*compute.Firewall, error) {
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

func (f *FirewallRuleDummyClient) CreateFirewallRule(project string, rule *compute.Firewall) error {
	for _, r := range f.Rules[project] {
		if r.Name == rule.Name {
			return fmt.Errorf("Rule already exists")
		}
	}

	f.Rules[project] = append(f.Rules[project], rule)
	return nil
}
func (f *FirewallRuleDummyClient) UpdateFirewallRule(project string, rule *compute.Firewall) error {
	//TODO
	return nil
}

func (f *FirewallRuleDummyClient) DeleteFirewallRule(project string, name string) error {
	rules := f.Rules[project]
	for index, rule := range rules {
		// Swap with last and drop last
		if rule.Name == name {
			rules[index] = rules[len(rules)-1]
			f.Rules[project] = rules[:len(rules)-1]
			return nil
		}
	}
	// TODO double check if deletion of non existing rule trigger error or nil
	return nil
}

/*
 * Lets check !
 *
 */
func TestCreateApplicationFirewallRules(t *testing.T) {
	manager, _ := NewFirewallRuleDummyClient()

	project := "dummy-project"
	serviceProject := "dummy-service-project"
	application := "dummy-application"
	var rules FirewallRuleList
	rule := FirewallRule{
		Rule:       compute.Firewall{Name: "remote", Network: "global/networks/default", Allowed: []*compute.FirewallAllowed{&compute.FirewallAllowed{Ports: []string{"22", "3389"}, IPProtocol: "TCP"}}},
		CustomName: "allow-tcp-22-3389",
	}
	rules = append(rules, rule)

	app := ApplicationRules{
		Project:        project,
		ServiceProject: serviceProject,
		Application:    application,
		Rules:          rules,
	}

	err := CreateApplicationFirewallRules(manager, app)
	if err != nil {
		t.Fatalf("Something wrong during rule creation. Got error %v\n", err)
	}

	// Check properties
	// TODO create helper func
	expected := fmt.Sprintf("%s-%s-%s", serviceProject, application, rule.CustomName)
	if manager.Rules[project][0].Name != expected {
		t.Errorf("Name don't match format. Got %s, expected %s\n", manager.Rules[project][0].Name, expected)

	}

	// Inster existing rule should trigger error
	err = CreateApplicationFirewallRules(manager, app)
	if err == nil {
		t.Errorf("Expected error during insert if rule already exists")
	}

}

func TestListApplicationFirewallRules(t *testing.T) {
	// Add dummy content
	manager, _ := NewFirewallRuleDummyClient()

	project := "kubernetes-host-project"
	serviceProjects := []string{"kubernetes-demo", "kubernetes-training"}
	applications := []string{"the-hard-way", "the-easy-way"}

	for _, serviceProject := range serviceProjects {
		for _, application := range applications {
			name := fmt.Sprintf("%s-%s-%s", serviceProject, application, "allow-external")
			rule := compute.Firewall{Name: name, Network: "global/networks/default", Allowed: []*compute.FirewallAllowed{&compute.FirewallAllowed{Ports: []string{"22", "6443"}, IPProtocol: "TCP"}}}
			manager.Rules[project] = append(manager.Rules[project], &rule)
		}
	}

	// Ask for non-existing project
	_, err := ListApplicationFirewallRules(manager, "non-existing-project", serviceProjects[0], applications[0])
	if err == nil {
		t.Fatalf("Expected error during ListApplicationFirewallRules on a non-existing project")
	}

	// Ask for one application in a random project
	rules, err := ListApplicationFirewallRules(manager, project, serviceProjects[0], applications[0])
	if err != nil {
		t.Fatalf("Something wrong during ListApplicationFirewallRules. Got error : %v\n", err)
	}

	if len(rules.Rules) != 1 {
		t.Errorf("Wrong rules count for project %s and application %s", serviceProjects[0], applications[0])
	}
}

func TestDeleteApplicationFirewallRules(t *testing.T) {
	// Add dummy content
	manager, _ := NewFirewallRuleDummyClient()

	project := "nginx-host-project"
	serviceProjects := []string{"nginx-demo", "nginx-training"}
	applications := []string{"front"}

	for _, serviceProject := range serviceProjects {
		for _, application := range applications {
			name := fmt.Sprintf("%s-%s-%s", serviceProject, application, "allow-publicly")
			rule := compute.Firewall{Name: name, Network: "global/networks/default", Allowed: []*compute.FirewallAllowed{&compute.FirewallAllowed{Ports: []string{"80", "443"}, IPProtocol: "TCP"}}}
			manager.Rules[project] = append(manager.Rules[project], &rule)
		}
	}

	fmt.Printf("%+v\n", manager.Rules)

	// Ask for delete application in nginx-demo project
	app := ApplicationRules{Project: project, ServiceProject: serviceProjects[0], Application: applications[0]}
	err := DeleteApplicationFirewallRules(manager, app)
	if err != nil {
		t.Fatalf("Unexpected error during Delete. Got %v\n", err)
	}

	// Verify empty rules
	apprules, err := ListApplicationFirewallRules(manager, project, serviceProjects[0], applications[0])
	if err != nil {
		t.Errorf("Unexpected error during ListApplicationFirewallRules for project %v", serviceProjects[0])
	}
	expected := 0
	if len(apprules.Rules) != expected {
		t.Errorf("Bad rules count for service project %v and application %v. Got %d expected %d", serviceProjects[1], applications[0], len(apprules.Rules), expected)
	}

	// Verify Others projects still exists
	apprules, err = ListApplicationFirewallRules(manager, project, serviceProjects[1], applications[0])
	if err != nil {
		t.Errorf("Unexpected error during ListApplicationFirewallRules for project %v", serviceProjects[1])
	}

	expected = 1
	if len(apprules.Rules) != expected {
		t.Errorf("Bad rules count for service project %v and application %v. Got %d expected %d", serviceProjects[1], applications[0], len(apprules.Rules), expected)
	}

	// Try to delete on non-existing project
	app = ApplicationRules{Project: "non-existing", ServiceProject: serviceProjects[0], Application: applications[0]}
	err = DeleteApplicationFirewallRules(manager, app)
	if err == nil {
		t.Fatalf("Expected error during Delete on non existing project. Got %v\n", err)
	}

}
