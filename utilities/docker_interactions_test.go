package utilities

import (
	"context"
	"testing"
)

var ctx = context.Background()
// const composeFilePath = "../docker-compose--for-test.yaml"
const composeFilePath = "/home/k/programming/caddy-proxy-organizer/docker-compose--for-test.yaml"
const serverContainerLabelsKey = "webserver-component"

func TestGetContainersAll(t *testing.T) {
	// -- Test if can capture any running containers

	// Make sure atleast one container is running
	runDockerComposeUp(ctx, RunDockerComposeParams{[]string{composeFilePath}, false})

	containers := GetContainersAll(ctx, nil)
	if len(containers) == 0 {
		t.Error("Result returned with no running containers")
	}
}

func TextGetContainersByLabel(t *testing.T) {
	// -- Test if can get any containers with the label

	// Start containers
	runDockerComposeUp(ctx, RunDockerComposeParams{[]string{composeFilePath}, false})

	labelKey := serverContainerLabelsKey
	labelArgs := []string{labelKey, ""}
	containers := GetContainersByLabel(ctx, nil, labelArgs...)
	if len(containers) == 0 {
		t.Errorf("Result returned with no running containers of label key %s", labelKey)
	}
	for _, container := range containers {
		for k,_ := range container.Labels {
			if k != labelKey {
				t.Errorf("Tried get containers with label key %s, but got a container with key %s", labelKey, k)
			}
		}
	}
}

func TestStopContainersByLabel(t *testing.T) {
	// -- Test if can stop containers by labelKey

	// Start containers
	runDockerComposeUp(ctx, RunDockerComposeParams{[]string{composeFilePath}, false})

	labelArgs := []string{serverContainerLabelsKey, ""}

	// Stop all "webserver-component" labeled containers
	StopContainersByLabel(ctx, nil, labelArgs...)

	containers := GetContainersByLabel(ctx, nil, labelArgs...)
	if len(containers) > 0 {
		t.Error("Result returned with no running containers")
	}
}