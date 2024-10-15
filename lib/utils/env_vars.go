package utils

import (
	"os"
)

type EnvVars struct {
	Environment Environment
	Domain      string
}

/* Load environment varaiables and return as a struct. */
func GetEnvVars() *EnvVars {
	var envVars EnvVars
	envVars.Environment = Environment(os.Getenv("ENVIRONMENT"))
	if envVars.Environment == "" {
		panic("failure to load the 'ENVIRONMENT' environment variable.")
	}
	envVars.Domain = os.Getenv("DOMAIN")
	if envVars.Domain == "" {
		panic("failure to load the 'DOMAIN' environment variable.")
	}

	return &envVars
}
