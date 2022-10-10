package utility

import (
	"bytes"
	"text/template"
)

type ServiceInfo struct {
	Name string
	Domain string
	Address string
	Ports []uint16
}

type ConsulTemplateImitator struct {
	Services []ServiceInfo
}

func (cti *ConsulTemplateImitator) ParseTemplate(filePath string) string {
	funcMap := template.FuncMap{
		"services": func() []ServiceInfo {
			return cti.Services
		},
	}
	
	templateJson, err := template.
		New(filePath).
		Funcs(funcMap).
		ParseFiles(filePath)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = templateJson.Execute(&buf, nil)
	if err != nil {
		ErrorLog.Fatalf("execution: %s", err)
	}

	return buf.String()
}