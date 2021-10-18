package usercss

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	ucPass = `/*==UserStyle==
@name         Name
@namespace    namespace
@description  Description
@author       Temp   <temp@example.com>		(https://temp.example.com)
@homepageURL  https://temp.example.com/temp
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
	ucFail = `/*==UserStyle==
@name
@namespace
@description  Description
@author       Temp <temp@example.com> (https://temp.example.com)
@homepageURL  https://temp.example.com/temp
@supportURL   https://temp.example.com/temp/issues
@updateURL    https://temp.example.com/temp/raw/temp.user.styl
@version
@license      MIT
@preprocessor uso
==/UserStyle== */

@-moz-document url(https://example.com/test) {
	:root {}
}

@-moz-document domain("example.com"), domain('example.org') {
	:root { --hello: 'world' }
}`
	domain = `/*==UserStyle==
@name         Name
@namespace    namespace
@description  Description
@author       Temp <temp@example.com> (https://temp.example.com)
@homepageURL  https://temp.example.com/temp
@supportURL   https://temp.example.com/temp/issues
@updateURL    https://temp.example.com/temp/raw/temp.user.styl
@version      1.0.0
@license      MIT
@preprocessor uso
==/UserStyle== */

@-moz-document domain("do as i say, not as i do"), regexp(^https?://(.+\.userstyles.world|localhost:[0-9]+)/[a-zA-Z0-9]{32,128}$) ,   domain(example.com)   , url-prefix("http://localhost")   {
	:root {}
}`
	tabs = `/*==UserStyle==
@name		newstyle
@namespace	somespace
@version	1.0.1
==/UserStyle== */`
	updateURL = "https://example.com/api/style/1.user.css"
)

func TestValidationPass(t *testing.T) {
	t.Parallel()

	uc := new(UserCSS)
	if err := uc.Parse(ucPass); err != nil {
		t.Fatal(err)
	}
	if err := uc.Validate(); err != nil {
		t.Fatal("Passing validation shouldn't return errors.")
	}
}

func TestValidationFail(t *testing.T) {
	t.Parallel()

	uc := new(UserCSS)
	if err := uc.Parse(ucFail); err != nil {
		t.Fatal(err)
	}
	if err := uc.Validate(); err == nil {
		t.Fatal("Failing validation should return errors.")
	}
}

func TestAuthor(t *testing.T) {
	t.Parallel()

	have := new(UserCSS)
	have.Parse(ucPass)
	want := UserCSS{
		Author: Author{
			Name:    "Temp",
			Email:   "temp@example.com",
			Website: "https://temp.example.com",
		},
	}

	if have.Author != want.Author {
		t.Fatal("Parsed author doesn't match.")
	}
}

func TestSingleDomain(t *testing.T) {
	t.Parallel()

	have := new(UserCSS)
	have.Parse(domain)
	want := UserCSS{
		MozDocument: []Domain{
			{
				Key:   "domain",
				Value: `do as i say, not as i do`,
			},
		},
	}

	if !reflect.DeepEqual(have.MozDocument[0], want.MozDocument[0]) {
		t.Fatal("Domains don't match.")
	}
}

func TestMultipleDomains(t *testing.T) {
	t.Parallel()

	have := new(UserCSS)
	have.Parse(ucPass)
	want := UserCSS{
		MozDocument: []Domain{
			{
				Key:   "url",
				Value: "https://example.com/test",
			},
			{
				Key:   "domain",
				Value: "example.com",
			},
			{
				Key:   "domain",
				Value: "example.org",
			},
		},
	}

	if !reflect.DeepEqual(have.MozDocument, want.MozDocument) {
		t.Fatal("Domains don't match.")
	}
}

func TestValidRemoteUserCSS(t *testing.T) {
	t.Parallel()

	URL := "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl"

	// Test will fail if URL is invalid.
	uc := new(UserCSS)
	if err := uc.ParseURL(URL); err != nil {
		t.Fatal(err)
	}

	if errs := uc.Validate(); errs != nil {
		t.Fatal(errs)
	}
}

func TestInvalidRemoteUserCSS(t *testing.T) {
	t.Parallel()

	URL := "https:///raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl"

	// Test will fail because protocol has three slashes instead of two.
	uc := new(UserCSS)
	if err := uc.ParseURL(URL); err == nil {
		t.Fatalf("Error parsing from URL: %v", err)
	}
}

func TestUserCSS(t *testing.T) {
	t.Parallel()

	have := new(UserCSS)
	have.Parse(ucPass)
	want := &UserCSS{
		Name:         "Name",
		Namespace:    "namespace",
		Description:  "Description",
		Version:      "1.0.0",
		License:      "MIT",
		HomepageURL:  "https://temp.example.com/temp",
		SupportURL:   "https://temp.example.com/temp/issues",
		UpdateURL:    "https://temp.example.com/temp/raw/temp.user.styl",
		Preprocessor: "uso",
		SourceCode:   fmt.Sprintf("%v", ucPass),
		Author: Author{
			Name:    "Temp",
			Email:   "temp@example.com",
			Website: "https://temp.example.com",
		},
		MozDocument: []Domain{
			{
				Key:   "url",
				Value: "https://example.com/test",
			},
			{
				Key:   "domain",
				Value: "example.com",
			},
			{
				Key:   "domain",
				Value: "example.org",
			},
		},
	}

	if !reflect.DeepEqual(have.MozDocument, want.MozDocument) {
		t.Fatal("UserCSS structs don't match.")
	}
}

func TestParseEmptyStyle(t *testing.T) {
	t.Parallel()

	uc := new(UserCSS)
	if err := uc.Parse(""); err == nil {
		t.Fatalf("Empty style should return an err, got: %v", err)
	}
}

func TestParseEmptyMetadata(t *testing.T) {
	t.Parallel()

	uc := new(UserCSS)
	if err := uc.Parse(ucPass); err != nil {
		t.Fatalf("Empty metadata should return an err, got: %v", err)
	}
}

func TestOverrideUpdateURL(t *testing.T) {
	t.Parallel()

	uc := new(UserCSS)
	if err := uc.Parse(ucPass); err != nil {
		t.Fatal(err)
	}

	uc.OverrideUpdateURL(updateURL)

	if uc.UpdateURL != updateURL {
		t.Fatal("Failed to override @updateURL field.")
	}
}

func TestFastOverrideUpdateURL(t *testing.T) {
	t.Parallel()

	have := new(UserCSS)
	code := OverrideUpdateURL(ucPass, updateURL)
	if err := have.Parse(code); err != nil {
		t.Fatal(err)
	}

	want := &UserCSS{UpdateURL: updateURL}

	if !reflect.DeepEqual(have.UpdateURL, want.UpdateURL) {
		t.Fatal("UpdateURL fields don't match.")
	}
}

func TestMetadataWithTabs(t *testing.T) {
	t.Parallel()

	have := new(UserCSS)
	have.Parse(tabs)
	want := &UserCSS{
		Name:       "newstyle",
		Namespace:  "somespace",
		Version:    "1.0.1",
		SourceCode: fmt.Sprintf("%v", tabs),
	}

	if !reflect.DeepEqual(have, want) {
		t.Fatal("UserCSS structs don't match.")
	}
}
