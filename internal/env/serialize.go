package env

import (
	"fmt"
	"strings"
)

// Serialize converts an env map into a .env file string.
// Keys are written in sorted order. Values containing spaces or special
// characters are double-quoted.
func Serialize(env map[string]string) string {
	keys := Keys(env)
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if needsQuoting(v) {
			fmt.Fprintf(&sb, "%s=\"%s\"\n", k, escapeValue(v))
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}

// needsQuoting returns true if the value should be wrapped in double quotes.
func needsQuoting(v string) bool {
	if v == "" {
		return false
	}
	for _, ch := range v {
		if ch == ' ' || ch == '\t' || ch == '#' || ch == '=' || ch == '\'' || ch == '\\' {
			return true
		}
	}
	return false
}

// escapeValue escapes double quotes and backslashes inside a quoted value.
func escapeValue(v string) string {
	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, "\"", "\\\"")
	return v
}
