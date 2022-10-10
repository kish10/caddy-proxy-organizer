package main

import (
	"context"
	"flag"

	"github.com/kish10/caddy-proxy-organizer/utility"
)

var fromDocker bool

func main() {

	flag.BoolVar(
		&fromDocker,
		"from-docker", 
		true, 
		"Specifies whether to generate caddy config using Docker container labels & API (default is true)",
	) 
	flag.Parse()

	if fromDocker {
		ctx := context.Background()
		cli := utility.GetDockerClient()
		LoadCaddyProxyJson(ctx, cli)
	} else {
		utility.InfoLog.Print("caddy-proxy-organizer ran but did nothing. Try using `--from-docker` flag.")
	}
}

