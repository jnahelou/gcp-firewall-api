package models

import (
	"encoding/json"
	"fmt"

	"google.golang.org/api/googleapi"
)

// GoogleApplicationError describe a custom error response
type GoogleApplicationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewGoogleApplicationError GoogleApplicationError constructor
func NewGoogleApplicationError(err *googleapi.Error) *GoogleApplicationError {
	return &GoogleApplicationError{
		Code:    err.Code,
		Message: err.Message,
	}
}

func (g *GoogleApplicationError) Error() string {
	return fmt.Sprintf("Error from Google: Code %d Message '%s'", g.Code, g.Message)
}

// JSON return error as JSON format
func (g *GoogleApplicationError) JSON() string {
	res, _ := json.Marshal(g)
	return string(res)
}
