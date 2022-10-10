package utility

import (
	"testing"
)

const testTemplateFilePath = "template--for-test.tmpl"

func TestParseTemplate(t *testing.T) {
	cti := ConsulTemplateImitator{
		Services: []ServiceInfo{
			ServiceInfo{
				Name: "test",
				Address: "127.0.0.0.1",
				Ports: [80],
			},
		},
	}

	// Test if compiles
	parsedTemplate := cti.ParseTemplate(testTemplateFilePath)
	if parsedTemplate == "" {
		t.Error("Returned empty result")
	}
}