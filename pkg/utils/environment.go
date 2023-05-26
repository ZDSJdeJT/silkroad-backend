package utils

import (
	"fmt"
	"os"
	"strings"
)

func CheckEnvVarsExist(vars []string) error {
	var missingVars []string
	for _, envVar := range vars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missingVars, ", "))
	}
	return nil
}
