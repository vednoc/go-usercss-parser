package usercss

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Error codes returned by parser failures.
var (
	ErrEmptyUserstyle = errors.New("userstyle is missing code")
	ErrEmptyMetadata  = errors.New("userstyle is missing metadata")
)

// Error codes returned by validation failures.
var (
	ErrEmptyName      = errors.New("@name field cannot be empty")
	ErrEmptyNamespace = errors.New("@namespace field cannot be empty")
	ErrEmptyVersion   = errors.New("@version field cannot be empty")
)

var (
	// Parse metadata.
	metaRe = regexp.MustCompile(`@.*`)

	// ParseAuthor regex.
	mailRe = regexp.MustCompile(`<(.*)>`)
	siteRe = regexp.MustCompile(`\((.*)\)`)

	// ParseDomains regex.
	areaRe    = regexp.MustCompile(`(?mU)(domain|url|url-prefix|regexp)\(.*\)`)
	areaKeyRe = regexp.MustCompile(`^\w+`)
	areaValRe = regexp.MustCompile(`\((.*)\)`)

	// OverrideUpdateURL regex.
	// `m` flag lets ^ and $ match beginning and end of multi-line matches.
	updateRe = regexp.MustCompile(`(?m)^(@updateURL\s*)(.+)$`)
)

type UserCSS struct {
	Name         string
	Namespace    string
	Description  string
	Version      string
	License      string
	HomepageURL  string
	SupportURL   string
	UpdateURL    string
	Preprocessor string
	SourceCode   string
	Author       Author
	MozDocument  []Domain
}

type Author struct {
	Name    string
	Email   string
	Website string
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

func (uc *UserCSS) ParseURL(url string) error {
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return uc.Parse(string(body))
}

func (uc *UserCSS) Parse(code string) error {
	uc.SourceCode = code
	if len(uc.SourceCode) == 0 {
		return ErrEmptyUserstyle
	}

	matches := metaRe.FindAllStringSubmatch(code, -1)
	if len(matches) == 0 {
		return ErrEmptyMetadata
	}

	for _, match := range matches {
		for _, s := range match {
			parts := strings.Fields(s)

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
			case "@author":
				uc.ParseAuthor(tail)
			case "@-moz-document":
				tail = strings.TrimRight(tail, " {")
				uc.ParseDomains(tail)
			}
		}
	}

	return nil
}

func (uc *UserCSS) ParseAuthor(data string) {
	// Using strings.Fields will trim all whitespace.
	parts := strings.Fields(data)
	a := Author{}

	// Check if name is set.
	if len(parts) >= 1 {
		a.Name = parts[0]
	}

	// Check if e-mail is set.
	if len(parts) >= 2 {
		// This will return a slice of e-mails.
		s := mailRe.FindStringSubmatch(parts[1])

		if s != nil {
			// We want the second one.
			email := s[1]

			a.Email = email
		}
	}

	if len(parts) >= 3 {
		// This will return a slice of URLs.
		s := siteRe.FindStringSubmatch(parts[2])

		if s != nil {
			// We want the second one.
			ws := s[1]

			a.Website = ws
		}
	}

	uc.Author = a
}

func (uc *UserCSS) ParseDomains(data string) {
	parts := areaRe.FindAllString(data, -1)

	for _, v := range parts {
		if documentKeyword(v) {
			key := areaKeyRe.FindStringSubmatch(v)[0]
			val := areaValRe.FindStringSubmatch(v)[1]

			// Trim quotes.
			val = strings.Trim(val, "'\"")

			uc.MozDocument = append(uc.MozDocument, Domain{
				Key:   key,
				Value: val,
			})
		}
	}
}

func documentKeyword(key string) bool {
	keys := []string{"url", "url-prefix", "regexp", "domain"}
	for _, v := range keys {
		if strings.HasPrefix(key, v) {
			return true
		}
	}

	return false
}

func (uc *UserCSS) Validate() Errors {
	var errors Errors

	if len(uc.Name) == 0 {
		errors = append(errors, Error{
			Name: "name",
			Code: ErrEmptyName,
		})
	}
	if len(uc.Namespace) == 0 {
		errors = append(errors, Error{
			Name: "namespace",
			Code: ErrEmptyNamespace,
		})
	}
	if len(uc.Version) == 0 {
		errors = append(errors, Error{
			Name: "version",
			Code: ErrEmptyVersion,
		})
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func (uc *UserCSS) OverrideUpdateURL(url string) {
	if uc.UpdateURL != "" {
		uc.UpdateURL = url

		// `${1}` will preserve whitespace between capturing group and the URL.
		uc.SourceCode = updateRe.ReplaceAllString(uc.SourceCode, "${1}"+url)
	}
}
