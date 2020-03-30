package helpers

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(ioutil.Discard)
}

type TestCase struct {
	Parameter     string
	ExpectedValue string
}

var results map[string]string = make(map[string]string)

func testFunction(r *http.Request) {
	project, serviceProject, application, rule := GetMuxVars(r)
	results["project"] = project
	results["service_project"] = serviceProject
	results["application"] = application
	results["rule"] = rule
}

func TestGetMuxVars(t *testing.T) {
	r := *&http.Request{}

	// Prepare test
	cases := []TestCase{
		TestCase{
			Parameter:     "project",
			ExpectedValue: "dummy_project",
		},
		TestCase{
			Parameter:     "service_project",
			ExpectedValue: "dummy_service_project",
		},
		TestCase{
			Parameter:     "application",
			ExpectedValue: "dummy_application",
		},
		TestCase{
			Parameter:     "rule",
			ExpectedValue: "dummy_roule",
		},
	}

	// Cast parameters to be understand by mux
	expected := map[string]string{}
	for _, c := range cases {
		expected[c.Parameter] = c.ExpectedValue
	}

	// Set parametters in the dummy request
	r = *mux.SetURLVars(&r, expected)

	// Execute tested function
	testFunction(&r)

	// Execute each tests
	for _, c := range cases {
		t.Run(c.Parameter, func(t *testing.T) {
			if results[c.Parameter] != c.ExpectedValue {
				t.Errorf("Expected %s got %s", c.ExpectedValue, results[c.Parameter])
			}
		})
	}

}
