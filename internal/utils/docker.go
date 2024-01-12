package utils

import "os"

func IsDocker() bool {
	v, exists := os.LookupEnv("OS_ENV")
	return exists && v == "docker"
}
