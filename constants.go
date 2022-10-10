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

// -- container labels

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

func LabelValueForServiceContainer() string {
	return getEnvVarWithDefault(
		"LABELS_VALUE_FOR_SERVICE_CONTAINER",
		"service",
	)
}

func LabelKeyForServiceDomain() string {
	return getEnvVarWithDefault(
		"LABELS_KEY_FOR_SERVICE_DOMAIN",
		"webserver-service-domain",
	)
}

// -- container network names

func NetworkNameForCaddyProxyExternal() string {
	return getEnvVarWithDefault(
		"NETWORK_NAME_FOR_CADDY_PROXY_EXTERNAL",
		"caddy-proxy-external-network",
	)
}


// -- necessary file paths

func pathCaddyProxyConfigJsonTemplate() string {
	return "caddy-config-template.json.tmpl"
}

func pathCaddyProxyConfigJson() string {
	return getEnvVarWithDefault(
		"PATH_CADDY_PROXY_JSON_CONFIG",
		"/usr/data/caddy_proxy_config/caddy.json",
	)
}
