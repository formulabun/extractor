package env

import "os"

var Host string

var Port string

func init() {
	portEnv := defaulted(os.Getenv("extractor_port"), "8080")
	hostEnv := defaulted(os.Getenv("extractor_host"), "extractor")

	Port = ":" + portEnv
	Host = hostEnv + Port
}

func defaulted(value, def string) string {
	if value == "" {
		return def
	}
	return value
}
