package main

import (
	"fmt"
)

func ApiCaddyLoadCmd(caddyJsonFilePath string) []string {
	localCurlCallCmd := []string{
		"curl", "http://localhost:2019/load",
		"-H", "Content-Type: application/json",
		"-d", fmt.Sprintf("@%s", caddyJsonFilePath),
	}

	return localCurlCallCmd
}