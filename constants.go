package main

import (
	"os"

	// "github.com/kish10/caddy-proxy-organizer/utility"
)

func getEnvVarWithDefault(name string, _default string) string {
	value := os.Getenv(name)
	if value == "" {
		value = _default
	}

	return value
}

func LabelKeyForServerContainers() string {
	return getEnvVarWithDefault(
		"LABELS_KEY_FOR_SERVER_CONTAINER",
		"webserver-component",
	)
}

func LabelValueForCaddyProxyContainer() string {
	return getEnvVarWithDefault(
		"LABELS_VALUE_FOR_CADDY_PROXY_CONTAINER",
		"caddy-proxy",
	)
}

func LabelValueForCaddyProxyOrganizerContainer() string {
	return getEnvVarWithDefault(
		"LABELS_VALUE_FOR_CADDY_PROXY_ORGANIZER_CONTAINER",
		"caddy-proxy-organizer",
	)
}

