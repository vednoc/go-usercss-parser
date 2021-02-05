package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	url  = "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl"
	temp = `/*==UserStyle==
@name         Name
@namespace    namespace
@description  Description
@author       Temp <temp@example.com> (https://temp.example.com)
@homepageURL  https://temp.example.com/temp/
@supportURL   https://temp.example.com/temp/issues
@updateURL    https://temp.example.com/temp/raw/temp.user.styl
@version      1.0.0
@license      MIT
@preprocessor uso
==/UserStyle== */

@-moz-document domain('example.com') {
	:root { --hello: 'world' }
}`
)

func ParseFromURL(url string) {
	req, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}

	Parse(string(body))
}

func Parse(data string) {
	r := regexp.MustCompile(`@.*`)
	matches := r.FindAllStringSubmatch(data, -1)

	// TODO: Store the data in a proper data structure.
	for _, match := range matches {
		for _, s := range match {
			parts := strings.Split(s, " ")
			head := parts[0]

			switch head {
			case "@name",
				"@namespace",
				"@description",
				"@author",
				"@version",
				"@license",
				"@homepageURL",
				"@supportURL",
				"@preprocessor",
				"@-moz-document":
				tail := strings.TrimSpace(strings.Join(parts[1:], " "))
				fmt.Printf("%-20s %s\n", head, tail)

				// TODO: Add the default case.
				// default:
				// 	fmt.Println("Not implemented yet!")
			}
		}
	}
}

func main() {
	Parse(temp)
	ParseFromURL(url)
}
