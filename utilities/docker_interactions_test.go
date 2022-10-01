package utilities

// github.com/kish10/caddy_proxy_organizer/

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestGetContainersAll(t *testing.T) {
	// -- Test if can capture any running containers

	// Make sure atleast one container is running
	runDockerComposeUp("docker-compose--for-test.yaml")

	containers := GetContainersAll(ctx, nil)
	if len(containers) == 0 {
		t.Error("Result returned with no running containers")
	}
}

func TextGetContainersByLabel(t *testing.T) {
	// -- Test if can get any containers with the label

	// Start containers
	runDockerComposeUp("docker-compose--for-test.yaml")

	labelKey := "webserver-component"
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
	runDockerComposeUp("docker-compose--for-test.yaml")

	labelArgs := []string{"webserver-component", ""}

	// Stop all "webserver-component" labeled containers
	StopContainersByLabel(ctx, nil, labelArgs...)

	containers := GetContainersByLabel(ctx, nil, labelArgs...)
	if len(containers) > 0 {
		t.Error("Result returned with no running containers")
	}
}