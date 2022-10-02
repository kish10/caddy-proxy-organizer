package utilities

import (
	"context"
	"testing"
)

var ctx = context.Background()
// const composeFilePath = "../docker-compose--for-test.yaml"
const composeFilePath = "/home/k/programming/caddy-proxy-organizer/docker-compose--for-test.yaml"
const labelKeyServerContainer = "webserver-component"
const labelValueCaddyProxy = "caddy-proxy"

func TestGetContainersAll(t *testing.T) {
	// -- Test if can capture any running containers

	// Make sure atleast one container is running
	RunDockerComposeUp(ctx, RunDockerComposeParams{[]string{composeFilePath}, false})

	containers := GetContainersAll(ctx, nil)
	if len(containers) == 0 {
		t.Error("Result returned with no running containers")
	}
}

func TestGetContainersByLabel(t *testing.T) {
	// -- Test if can get any containers with the label

	// Start containers
	RunDockerComposeUp(ctx, RunDockerComposeParams{[]string{composeFilePath}, false})

	cases := map[int][]string{
		1: []string{labelKeyServerContainer, ""},
		2: []string{"", labelValueCaddyProxy},
		3: []string{labelKeyServerContainer, labelValueCaddyProxy},
	}

	for caseNum, labelKeyValue := range cases {
		containers := GetContainersByLabel(ctx, nil, labelKeyValue)

		if len(containers) == 0 {
			t.Errorf("Result returned with no running containers of label key `%s` & value `%s`", labelKeyValue[0], labelKeyValue[1])
		}

		for _, container := range containers {
			keyFound := false
			valueFound := false
			for k,v := range container.Labels {
				
				if k == labelKeyValue[0] {
					keyFound = true
				}

				if v == labelKeyValue[1] {
					valueFound = true
				}
			}

			testFail := false
			switch {
			case caseNum == 1 && !keyFound:
				testFail = true
			case caseNum == 2 && !valueFound:
				testFail = true
			case caseNum == 3 && !keyFound && !valueFound:
				testFail = true
			}
	
			if testFail {
				t.Errorf("Tried get containers with label key `%s` & value `%s`, but got a container without any of those", labelKeyValue[0], labelKeyValue[1])
			}
		}
	}
}

func TestStopContainersByLabel(t *testing.T) {
	// -- Test if can stop containers by labelKey

	// Start containers
	RunDockerComposeUp(ctx, RunDockerComposeParams{[]string{composeFilePath}, false})

	labelKeyValue := []string{labelKeyServerContainer, ""}

	// Stop all "webserver-component" labeled containers
	StopContainersByLabel(ctx, nil, labelKeyValue)

	containers := GetContainersByLabel(ctx, nil, labelKeyValue)
	if len(containers) > 0 {
		t.Error("Result returned with no running containers")
	}
}