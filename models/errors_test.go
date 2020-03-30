package models

import (
	"fmt"
	"testing"

	"google.golang.org/api/googleapi"
)

type TestCase struct {
	Title    string
	Expected interface{}
	Got      interface{}
}

func TestGoogleApplicationError(t *testing.T) {
	expectedError := *&googleapi.Error{
		Code:    1234,
		Message: "dummy",
	}

	// Execute function
	testedError := NewGoogleApplicationError(&expectedError)

	suite := []TestCase{
		TestCase{
			Title:    "Error should not be nil",
			Expected: false,
			Got:      testedError == nil,
		},
		TestCase{
			Title:    "Error code should be identical",
			Expected: expectedError.Code,
			Got:      testedError.Code,
		},
		TestCase{
			Title:    "Error messaage shoud be identical",
			Expected: expectedError.Message,
			Got:      testedError.Message,
		},
		TestCase{
			Title:    "Error() method should return message as fromatted string",
			Expected: fmt.Sprintf("Error from Google: Code %d Message '%s'", expectedError.Code, expectedError.Message),
			Got:      testedError.Error(),
		},
		TestCase{
			Title:    "JSON() method should return error as JSON",
			Expected: `{"code":1234,"message":"dummy"}`,
			Got:      testedError.JSON(),
		},
	}

	// Launch test
	for _, suiteCase := range suite {
		t.Run(suiteCase.Title, func(t *testing.T) {
			if suiteCase.Expected != suiteCase.Got {
				t.Errorf("Got %d want %d", suiteCase.Got, suiteCase.Expected)
			}
		})
	}
}
