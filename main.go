package main

import (
	"net/http"
	"os"

	"github.com/dgaldamez77/loc/endpoints"
	"github.com/dgaldamez77/loc/log"
)

func main() {
	validateEnvironment()

	router := endpoints.RegisterEndpoints()

	port := "8080"
	log.LogInfo("Service listening on port " + port)

	// for a real implementation, this needs to be changed to a serve function
	// that suports time out
	//nolint
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.LogError(err, "failed to start server")
	}
}

func validateEnvironment() {
	envVars := []string{"DB_URI"}
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			panic("unable to find environment variable " + envVar)
		}
	}
}
