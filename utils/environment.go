package utils

import (
	"fmt"
	"os"
	"strings"
)

const (
	AppMode          = "APP_MODE"
	APPPort          = "APP_PORT"
	APPName          = "APP_NAME"
	APPVersion       = "APP_VERSION"
	DatabaseDSN      = "DATABASE_DSN"
	JWTSecretKey     = "JWT_SECRET_KEY"
	JWTExpireMinutes = "JWT_EXPIRE_MINUTES"
)

var initialEnvVars = []string{
	APPPort,
	APPName,
	APPVersion,
	DatabaseDSN,
	JWTSecretKey,
	JWTExpireMinutes,
}

func CheckEnvVarsExist() error {
	var missingVars []string
	for _, envVar := range initialEnvVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missingVars, ", "))
	}
	return nil
}
