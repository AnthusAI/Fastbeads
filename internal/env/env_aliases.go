package env

import "os"

// DefaultEnvPrefixes defines the env var prefixes checked by LookupEnvAlias.
// Order matters: earlier prefixes take priority when multiple are set.
var DefaultEnvPrefixes = []string{"FBD", "BD", "BEADS"}

// LookupEnvAlias checks for an environment variable using common prefixes.
// Example: base="ACTOR" checks FBD_ACTOR, then BD_ACTOR, then BEADS_ACTOR.
func LookupEnvAlias(base string, prefixes ...string) (string, bool) {
	if base == "" {
		return "", false
	}
	if len(prefixes) == 0 {
		prefixes = DefaultEnvPrefixes
	}
	for _, prefix := range prefixes {
		key := prefix + "_" + base
		if val, ok := os.LookupEnv(key); ok && val != "" {
			return val, true
		}
	}
	return "", false
}

// GetEnvAlias returns the first non-empty value found by LookupEnvAlias.
func GetEnvAlias(base string, prefixes ...string) string {
	if val, ok := LookupEnvAlias(base, prefixes...); ok {
		return val
	}
	return ""
}
