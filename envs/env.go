package envs

import (
	"os"
)

func FindEnv(key, defaultValue string) string {
	value, defined := os.LookupEnv(key)
	if defined {
		return value
	}
	return defaultValue
}
