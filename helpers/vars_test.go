package helpers

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

type testCase struct {
	Parameter     string
	ExpectedValue string
	TestFuncton   func(r *http.Request) string
}

func Project(r *http.Request) string {
	v, _, _, _ := GetMuxVars(r)
	return v
}

func ServiceProject(r *http.Request) string {
	_, v, _, _ := GetMuxVars(r)
	return v
}

func Application(r *http.Request) string {
	_, _, v, _ := GetMuxVars(r)
	return v
}

func Rule(r *http.Request) string {
	_, _, _, v := GetMuxVars(r)
	return v
}

func TestGetMuxVars(t *testing.T) {
	r := *&http.Request{}

	// Prepare test
	cases := []testCase{
		testCase{
			Parameter:     "project",
			ExpectedValue: "dummy_project",
			TestFuncton:   Project,
		},
		testCase{
			Parameter:     "service-project",
			ExpectedValue: "dummy_service_project",
			TestFuncton:   ServiceProject,
		},
		testCase{
			Parameter:     "application",
			ExpectedValue: "dummy_application",
			TestFuncton:   Application,
		},
		testCase{
			Parameter:     "rule",
			ExpectedValue: "dummy_roule",
			TestFuncton:   Rule,
		},
	}

	// Cast parameters to be understand by mux
	expected := map[string]string{}
	for _, c := range cases {
		expected[c.Parameter] = c.ExpectedValue
	}

	// Set parametters in the dummy request
	r = *mux.SetURLVars(&r, expected)

	// Execute each tests
	for _, c := range cases {
		t.Run(c.Parameter, func(t *testing.T) {
			if result := c.TestFuncton(&r); result != c.ExpectedValue {
				t.Errorf("Expected %s got %s", c.ExpectedValue, result)
			}
		})
	}

}
