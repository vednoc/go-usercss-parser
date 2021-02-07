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

@-moz-document domain('example.com'), domain('example.org') {
	:root { --hello: 'world' }
}`
)

type UserCSS struct {
	Name         string
	Namespace    string
	Description  string
	Author       string
	Version      string
	License      string
	HomepageURL  string
	SupportURL   string
	UpdateURL    string
	Preprocessor string
	MozDocument  []string
}

func ParseFromURL(url string) *UserCSS {
	req, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}

	return Parse(string(body))
}

func Parse(data string) *UserCSS {
	r := regexp.MustCompile(`@.*`)
	matches := r.FindAllStringSubmatch(data, -1)

	uc := new(UserCSS)

	for _, match := range matches {
		for _, s := range match {
			parts := strings.Split(s, " ")

			// Metadata fields.
			head := parts[0]
			tail := strings.TrimSpace(strings.Join(parts[1:], " "))

			switch head {
			case "@name":
				uc.Name = tail
			case "@namespace":
				uc.Namespace = tail
			case "@description":
				uc.Description = tail
			case "@author":
				uc.Author = tail
			case "@version":
				uc.Version = tail
			case "@license":
				uc.License = tail
			case "@homepageURL":
				uc.HomepageURL = tail
			case "@supportURL":
				uc.SupportURL = tail
			case "@updateURL":
				uc.UpdateURL = tail
			case "@preprocessor":
				uc.Preprocessor = tail
			case "@-moz-document":
				tail = strings.TrimRight(tail, " {")
				uc.MozDocument = append(uc.MozDocument, tail)

				// TODO: Add the default case.
				// default:
				// 	fmt.Println("Not implemented yet!")
			}
		}
	}

	return uc
}

func main() {
	fmt.Printf("Parse temp data:\n%#+v\n", Parse(temp))
	fmt.Printf("Parse real data:\n%#+v\n", ParseFromURL(url))
}
