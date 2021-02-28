package curlparser

import (
	"io/ioutil"
	"log"
	"regexp"
)

type CurlParser struct {
	Headers map[string]string
	source  string
}

func (cp *CurlParser) GetConf() *CurlParser {
	file, err := ioutil.ReadFile("./curl.txt")
	if err != nil {
		log.Printf("Cant get config #%v", err)
	}

	str := string(file)

	cp.source = str
	cp.Headers = make(map[string]string)

	r := regexp.MustCompile(`(?m)-H '(.+): (.+)'`)

	for _, match := range r.FindAllStringSubmatch(str, -1) {
		cp.Headers[match[1]] = match[2]
	}

	return cp
}

func (cp *CurlParser) GetHeader(name string) string {
	value, ok := cp.Headers[name]

	if !ok {
		return ""
	}

	return value
}
