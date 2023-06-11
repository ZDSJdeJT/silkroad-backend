package utils

import (
	"fmt"
	"os"
	"strings"
)

const (
	AppMode          = "APP_MODE"
	JWTSecretKey     = "JWT_SECRET_KEY"
	JWTExpireMinutes = "JWT_EXPIRE_MINUTES"
)

var initialEnvVars = []string{
	JWTSecretKey,
	JWTExpireMinutes,
}

const (
	APPName     = "Silk Road"
	APPVersion  = "1.0.0"
	APPPort     = "4000"
	DatabaseDSN = "./data/database/sqlite.db"
)

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
