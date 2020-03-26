package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jnahelou/gcp-firewall-api/handlers"
	"github.com/jnahelou/gcp-firewall-api/helpers"
	"github.com/sirupsen/logrus"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Ok")
}

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
	managerRouter := r.PathPrefix("/project/{project}/service-project/{service-project}/application/{application}").Subrouter()
	managerRouter.Path("").Methods("GET").HandlerFunc(handlers.ListFirewallRulesHandler)
	managerRouter.Path("").Methods("POST").HandlerFunc(handlers.CreateFirewallRulesHandler)
	managerRouter.Path("").Methods("DELETE").HandlerFunc(handlers.DeleteFirewallRulesHandler)

	// Manage a specific rule
	ruleRouter := r.PathPrefix("/project/{project}/service-project/{service-project}/application/{application}/name/{rule}").Subrouter()
	ruleRouter.Path("").Methods("PUT").HandlerFunc(handlers.UpdateFirewallRuleHandler)
	ruleRouter.Path("").Methods("DELETE").HandlerFunc(handlers.DeleteFirewallRuleHandler)

	r.Path("/_health").Methods("GET").HandlerFunc(handleHealth)

	helpers.InitLogger()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	logrus.Printf("Listening on port %s", port)
	logrus.Print(srv.ListenAndServe())
}
