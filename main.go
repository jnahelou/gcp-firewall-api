package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adeo/iwc-gcp-firewall-api/handlers"
	"github.com/adeo/iwc-gcp-firewall-api/helpers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// log access log
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"request_uri": r.RequestURI,
			"user_agent":  r.UserAgent(),
		}).Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// define JSON as default return content type
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter().StrictSlash(true)
	// Disable http access log on testing
	if os.Getenv("CI") == "" {
		r.Use(loggingMiddleware)
	}
	r.Use(contentTypeMiddleware)

	// Manage sets of rules
	managerRouter := r.PathPrefix("/project/{project}/service_project/{service_project}/application/{application}").Subrouter()
	managerRouter.Path("").Methods("GET").HandlerFunc(handlers.ListFirewallRuleHandler)

	// Manage a specific rule
	ruleRouter := r.PathPrefix("/project/{project}/service_project/{service_project}/application/{application}/firewall_rule/{rule}").Subrouter()
	ruleRouter.Path("").Methods("POST").HandlerFunc(handlers.CreateFirewallRuleHandler)
	ruleRouter.Path("").Methods("GET").HandlerFunc(handlers.GetFirewallRuleHandler)
	ruleRouter.Path("").Methods("DELETE").HandlerFunc(handlers.DeleteFirewallRuleHandler)

	r.Path("/_health").Methods("GET").HandlerFunc(handlers.HealthCheckHandler)

	helpers.InitLogger()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	logrus.Printf("Listening on port %s", port)
	logrus.Print(srv.ListenAndServe())
}
