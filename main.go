package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type Credentials struct {
	Login    string
	Password string
}

func main() {
	value, err := GenXMLPasswordLine("COOKIE", "CAKE")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", value)
}

func GenXMLPasswordLine(login string, password string) (string, error) {
	template_string := "<value><struct><member><name>methodName</name><value><string>wp.getUsersBlogs</string></value></member><member><name>params</name><value><array><data><value><array><data><value><string>{{.Login}}</string></value><value><string>{{.Password}}</string></value></data></array></value></data></array></value></member></struct></value>"
	var buf bytes.Buffer

	credentials := Credentials{login, password}

	template, err := template.New("xml_temlate").Parse(template_string)
	if err != nil {
		return "", fmt.Errorf("Failed to parse xml template: %w", err)
	}

	err = template.Execute(&buf, credentials)
	if err != nil {
		return "", fmt.Errorf("Failed to execute xml template: %w", err)
	}

	return buf.String(), nil

}
