package helpers

import (
	"errors"
	"net/http"
	"strings"
)

func getPathContent(r *http.Request) (string, string, error) {
	urlPart := strings.Split(r.URL.Path, "/")
	if len(urlPart) != 6 {
		return "", "", errors.New("Unexpected path content")
	}

	project := urlPart[2]
	application := urlPart[4]

	//TODO verify content
	return project, application, nil
}
