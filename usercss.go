package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	// Errors.
	ErrEmptyName      = errors.New("name cannot be empty")
	ErrEmptyNamespace = errors.New("namespace cannot be empty")
	ErrEmptyVersion   = errors.New("version cannot be empty")

	// Test data.
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

@-moz-document url(https://example.com/test) {
	:root {}
}

@-moz-document domain("example.com"), domain('example.org') {
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
	MozDocument  []Domain
}

type Domain struct {
	Key   string
	Value string
}

type Error struct {
	Name string
	Code error
}

type Errors []Error

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

	return ParseFromString(string(body))
}

func ParseFromString(data string) *UserCSS {
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
				ParseDomains(tail, uc)
			}
		}
	}

	return uc
}

func ParseDomains(data string, uc *UserCSS) {
	parts := strings.Split(data, ",")

	// Regex rules.
	kr := regexp.MustCompile(`^\w+`)
	vr := regexp.MustCompile(`\((.*)\)`)

	for _, v := range parts {
		trim := strings.TrimSpace(v)
		key := kr.FindStringSubmatch(trim)[0]
		val := vr.FindStringSubmatch(trim)[1]

		// Trim quotes.
		val = strings.Trim(val, "'\"")

		uc.MozDocument = append(uc.MozDocument, Domain{
			Key:   key,
			Value: val,
		})
	}
}

func BasicMetadataValidation(uc *UserCSS) (bool, Errors) {
	errors := Errors{}

	if len(uc.Name) == 0 {
		err := Error{Name: "domain", Code: ErrEmptyName}
		errors = append(errors, err)
	}
	if len(uc.Namespace) == 0 {
		err := Error{Name: "namespace", Code: ErrEmptyNamespace}
		errors = append(errors, err)
	}
	if len(uc.Version) == 0 {
		err := Error{Name: "version", Code: ErrEmptyVersion}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}

func main() {
	temp := ParseFromString(temp)
	real := ParseFromURL(url)

	fmt.Printf("Temp data:\n%#+v\n", temp)
	fmt.Printf("Real data:\n%#+v\n", real)

	validateTemp, err := BasicMetadataValidation(temp)
	fmt.Printf("Temp data validation: %v\n", validateTemp)
	if validateTemp == false {
		for name, code := range err {
			fmt.Println("Error:", name, code)
		}
	}

	validateReal, err := BasicMetadataValidation(real)
	fmt.Printf("Real data validation: %v\n", validateReal)
	if validateReal == false {
		for name, code := range err {
			fmt.Println("Error:", name, code)
		}
	}
}
