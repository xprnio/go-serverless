package utils

import "fmt"

func EnvMapToList(m map[string]string) []string {
	var env []string

	for key, val := range m {
		item := fmt.Sprintf("%s=%s", key, val)
		env = append(env, item)
	}

	return env
}
