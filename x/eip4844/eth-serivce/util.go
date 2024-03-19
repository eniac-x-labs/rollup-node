package eth_serivce

import (
	"os"
	"reflect"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"
)

// PrefixEnvVar adds a prefix to the environment variable,
// and returns the env-var wrapped in a slice for usage with urfave CLI v2.
func PrefixEnvVar(prefix, suffix string) []string {
	return []string{prefix + "_" + suffix}
}

// ValidateEnvVars logs all env vars that are found where the env var is
// prefixed with the supplied prefix (like OP_BATCHER) but there is no
// actual env var with that name.
// It helps validate that the supplied env vars are in fact valid.
func ValidateEnvVars(prefix string, flags []cli.Flag, log log.Logger) {
	for _, envVar := range validateEnvVars(prefix, os.Environ(), cliFlagsToEnvVars(flags)) {
		log.Warn("Unknown env var", "prefix", prefix, "env_var", envVar)
	}
}

func cliFlagsToEnvVars(flags []cli.Flag) map[string]struct{} {
	definedEnvVars := make(map[string]struct{})
	for _, flag := range flags {
		envVars := reflect.ValueOf(flag).Elem().FieldByName("EnvVars")
		for i := 0; i < envVars.Len(); i++ {
			envVarField := envVars.Index(i)
			definedEnvVars[envVarField.String()] = struct{}{}
		}
	}
	return definedEnvVars
}

// validateEnvVars returns a list of the unknown environment variables that match the prefix.
func validateEnvVars(prefix string, providedEnvVars []string, definedEnvVars map[string]struct{}) []string {
	var out []string
	for _, envVar := range providedEnvVars {
		parts := strings.Split(envVar, "=")
		if len(parts) == 0 {
			continue
		}
		key := parts[0]
		if strings.HasPrefix(key, prefix) {
			if _, ok := definedEnvVars[key]; !ok {
				out = append(out, envVar)
			}
		}
	}
	return out
}
