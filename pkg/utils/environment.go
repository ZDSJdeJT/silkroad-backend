package utils

import (
	"fmt"
	"os"
)

func CheckEnvVarsExist(vars []string) (error, error) {
	for _, envVar := range vars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("missing environment variable %s", envVar), nil
		}
	}
	return nil, nil
}
