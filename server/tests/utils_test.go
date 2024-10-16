package tests

import (
	"os"
	"testing"

	"xes.software/rpg/lib/utils"
)

func TestEnvVar(t *testing.T) {
	os.Setenv("ENVIRONMENT", string(utils.Production))
	os.Setenv("DOMAIN", "rpg.xes.software")

	envVars := utils.GetEnvVars()
	if envVars.Domain != os.Getenv("DOMAIN") {
		t.Errorf("envVars.Domain:%s != os.Getenv(\"DOMAIN\"):%s", envVars.Domain, os.Getenv("DOMAIN"))
	}
	if envVars.Environment != utils.Environment(os.Getenv("ENVIRONMENT")) {
		t.Errorf("envVars.Environment:%s != os.Getenv(\"ENVIRONMENT\"):%s", envVars.Environment, os.Getenv("ENVIRONMENT"))
	}
}
