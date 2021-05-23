package usercss

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	// Validation errors.
	ErrEmptyName         = errors.New("@name field cannot be empty")
	ErrEmptyNamespace    = errors.New("@namespace field cannot be empty")
	ErrEmptyVersion      = errors.New("@version field cannot be empty")
	ErrMultipleOccurence = errors.New("duplicate occurence has been found in metadata")
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
	HintErrors   Errors
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

func ParseFromURL(url string) (*UserCSS, error) {
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return ParseFromString(string(body)), nil
}

func ParseFromString(data string) *UserCSS {
	r := regexp.MustCompile(`@.*`)
	matches := r.FindAllStringSubmatch(data, -1)

	uc := &UserCSS{
		SourceCode: data,
	}

	for _, match := range matches {
		for _, s := range match {
			parts := strings.Fields(s)

			// Metadata fields.
			head := parts[0]
			tail := strings.TrimSpace(strings.Join(parts[1:], " "))

			switch head {
			case "@name":
				if uc.Name != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "name", Code: ErrMultipleOccurence})
				}
				uc.Name = tail
			case "@namespace":
				if uc.Namespace != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "namespace", Code: ErrMultipleOccurence})
				}
				uc.Namespace = tail
			case "@description":
				if uc.Description != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "description", Code: ErrMultipleOccurence})
				}
				uc.Description = tail
			case "@version":
				if uc.Version != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "version", Code: ErrMultipleOccurence})
				}
				uc.Version = tail
			case "@license":
				if uc.License != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "license", Code: ErrMultipleOccurence})
				}
				uc.License = tail
			case "@homepageURL":
				if uc.HomepageURL != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "license", Code: ErrMultipleOccurence})
				}
				uc.HomepageURL = tail
			case "@supportURL":
				if uc.SupportURL != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "license", Code: ErrMultipleOccurence})
				}
				uc.SupportURL = tail
			case "@updateURL":
				if uc.UpdateURL != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "license", Code: ErrMultipleOccurence})
				}
				uc.UpdateURL = tail
			case "@preprocessor":
				if uc.Preprocessor != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "license", Code: ErrMultipleOccurence})
				}
				uc.Preprocessor = tail
			case "@author":
				if uc.Author.Name != "" {
					uc.HintErrors = append(uc.HintErrors, Error{Name: "license", Code: ErrMultipleOccurence})
				}
				ParseAuthor(tail, uc)
			// Multiple @-moz-document is allowed and supported.
			case "@-moz-document":
				tail = strings.TrimRight(tail, " {")
				ParseDomains(tail, uc)
			}
		}
	}

	return uc
}

func ParseAuthor(data string, uc *UserCSS) {
	// Using strings.Fields will trim all whitespace.
	parts := strings.Fields(data)
	a := Author{}

	// Check if name is set.
	if len(parts) >= 1 {
		a.Name = parts[0]
	}

	// Check if e-mail is set.
	if len(parts) >= 2 {
		er := regexp.MustCompile(`<(.*)>`)

		// This will return a slice of e-mails.
		s := er.FindStringSubmatch(parts[1])

		if s != nil {
			// We want the second one.
			email := s[1]

			a.Email = email
		}
	}

	if len(parts) >= 3 {
		wr := regexp.MustCompile(`\((.*)\)`)

		// This will return a slice of URLs.
		s := wr.FindStringSubmatch(parts[2])

		if s != nil {
			// We want the second one.
			ws := s[1]

			a.Website = ws
		}
	}

	uc.Author = a
}

func ParseDomains(data string, uc *UserCSS) {
	re := regexp.MustCompile(`(?mU)(domain|url|url-prefix|regexp)\(.*\)`)
	parts := re.FindAllString(data, -1)

	// Regex rules.
	kr := regexp.MustCompile(`^\w+`)
	vr := regexp.MustCompile(`\((.*)\)`)

	for _, v := range parts {
		if documentKeyword(v) {
			key := kr.FindStringSubmatch(v)[0]
			val := vr.FindStringSubmatch(v)[1]

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

func BasicMetadataValidation(uc *UserCSS) Errors {
	errors := Errors{}

	if len(uc.HintErrors) != 0 {
		errors = append(errors, uc.HintErrors...)
	}
	if len(uc.Name) == 0 {
		err := Error{Name: "name", Code: ErrEmptyName}
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
		return errors
	}

	return nil
}

func (uc *UserCSS) OverrideUpdateURL(url string) {
	if uc.UpdateURL != "" {
		uc.UpdateURL = url

		// `m` flag will let ^ and $ match start/end of lines.
		re := regexp.MustCompile(`(?m)^(@updateURL\s*)(.+)$`)

		// `${1}` allows us to keep whitespace between capturing group and URL.
		uc.SourceCode = re.ReplaceAllString(uc.SourceCode, "${1}"+url)
	}
}
