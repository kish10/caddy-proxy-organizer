package main

import (
	"fmt"
)

func ApiCaddyLoadCmd(caddyJsonFilePath string, caddyJsonString string) []string {
	if caddyJsonFilePath == "" && caddyJsonString == "" {
		panic("Need to provide json config data either as a string or filepath")
	}

	dataArg := fmt.Sprintf("@%s", caddyJsonFilePath)
	if caddyJsonString != "" {
		dataArg = caddyJsonString
	}

	localCurlCallCmd := []string{
		"curl", "http://localhost:2019/load",
		"-H", "Content-Type: application/json",
		"-d", dataArg,
	}

	return localCurlCallCmd
}