package tests

import (
	"os"
	"testing"

	"xes.software/rpg/lib/utils"
)

func TestLoadEnvVars(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", string(utils.Production))
	if err != nil {
		t.Errorf("error setting ENVIRONMENT env var.")
	}
	os.Setenv("DOMAIN", "game.xes.software")
	envVars := utils.GetEnvVars()
	if envVars.Domain != os.Getenv("DOMAIN") {
		t.Errorf("Environment variable 'DOMAIN' was not as expected envVars:%s != os.Getenv:%s", envVars.Domain, os.Getenv("DOMAIN"))
	}
	if envVars.Environment != utils.Environment(os.Getenv("ENVIRONMENT")) {
		t.Errorf("Environment variable 'ENVIRONMENT' was not as expected envVars:%s != os.Getenv:%s", envVars.Environment, os.Getenv("ENVIRONMENT"))
	}
}

func TestLogger(t *testing.T) {
}
