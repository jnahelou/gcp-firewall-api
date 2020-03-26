package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jnahelou/gcp-firewall-api/handlers"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	r := mux.NewRouter()
	r.Path("/_healthz").Methods("GET").HandlerFunc(handleHealth)
	r.Path("/project/{project}/service-project/{service-project}/application/{application}").Methods("GET").HandlerFunc(handlers.ListFirewallRulesHandler)
	r.Path("/project/{project}/service-project/{service-project}/application/{application}").Methods("POST").HandlerFunc(handlers.CreateFirewallRulesHandler)
	r.Path("/project/{project}/service-project/{service-project}/application/{application}").Methods("DELETE").HandlerFunc(handlers.DeleteFirewallRulesHandler)
	r.Path("/project/{project}/service-project/{service-project}/application/{application}/name/{rule}").Methods("POST").HandlerFunc(handlers.CreateFirewallRuleHandler)
	r.Path("/project/{project}/service-project/{service-project}/application/{application}/name/{rule}").Methods("DELETE").HandlerFunc(handlers.DeleteFirewallRuleHandler)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	log.Print(srv.ListenAndServe())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Ok")
}
