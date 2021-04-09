package usercss

import (
	"fmt"
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

@-moz-document regexp(^https?://(.+\.userstyles.world|localhost:[0-9]+)/[a-zA-Z0-9]{32,128}$) ,   domain(example.com) {
	:root {}
}`
	tabs = `/*==UserStyle==
@name		newstyle
@namespace	somespace
@version	1.0.1
==/UserStyle== */`
)

func TestValidationPass(t *testing.T) {
	uc := ParseFromString(ucPass)
	err := BasicMetadataValidation(uc)
	if err != nil {
		t.Fatal("Passed validation shouldn't return errors.")
	}
}

func TestValidationFail(t *testing.T) {
	uc := ParseFromString(ucFail)
	err := BasicMetadataValidation(uc)
	if err == nil {
		t.Fatal("Failed validation should return errors.")
	}
}

func TestAuthor(t *testing.T) {
	data := ParseFromString(ucPass)
	pass := Author{
		Name:    "Temp",
		Email:   "temp@example.com",
		Website: "https://temp.example.com",
	}

	dataString := fmt.Sprintf("%#+v", data.Author)
	passString := fmt.Sprintf("%#+v", pass)

	if dataString != passString {
		t.Fatal("Parsed author doesn't match.")
	}
}

func TestSingleDomain(t *testing.T) {
	data := ParseFromString(domain)
	pass := Domain{
		Key:   "regexp",
		Value: `^https?://(.+\.userstyles.world|localhost:[0-9]+)/[a-zA-Z0-9]{32,128}$`,
	}

	if data.MozDocument[0] != pass {
		t.Fatal("Domains don't match.")
	}
}

func TestMultipleDomains(t *testing.T) {
	data := ParseFromString(ucPass)
	pass := []Domain{
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
	}

	dataString := fmt.Sprintf("%#+v", data.MozDocument)
	passString := fmt.Sprintf("%#+v", pass)

	if dataString != passString {
		t.Fatal("Domain slices don't match.")
	}
}

func TestValidRemoteUserCSS(t *testing.T) {
	URL := "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl"

	// Test will fail if URL is invalid.
	data, err := ParseFromURL(URL)
	if err != nil {
		t.Fatal(err)
	}

	errs := BasicMetadataValidation(data)
	if errs != nil {
		t.Fatal(errs)
	}
}

func TestInvalidRemoteUserCSS(t *testing.T) {
	URL := "https:///raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl"

	// Test will fail because protocol has three slashes instead of two.
	_, err := ParseFromURL(URL)
	if err == nil {
		t.Fatalf("Error parsing from URL: %v", err)
	}
}

func TestUserCSS(t *testing.T) {
	data := ParseFromString(ucPass)
	pass := &UserCSS{
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

	dataString := fmt.Sprintf("%#+v", data)
	passString := fmt.Sprintf("%#+v", pass)

	if dataString != passString {
		t.Fatal("UserCSS structs don't match.")
	}
}

func TestOverrideUpdateURL(t *testing.T) {
	data := ParseFromString(ucPass)

	url := "https://example.com/api/style/1.user.css"
	data.OverrideUpdateURL(url)

	if data.UpdateURL != url {
		t.Fatal("Failed to override @updateURL field.")
	}
}

func TestMetadatawithTabs(t *testing.T) {
	data := ParseFromString(tabs)
	pass := &UserCSS{
		Name:       "newstyle",
		Namespace:  "somespace",
		Version:    "1.0.1",
		SourceCode: fmt.Sprintf("%v", tabs),
	}

	dataString := fmt.Sprintf("%#+v", data)
	passString := fmt.Sprintf("%#+v", pass)

	if dataString != passString {
		t.Fatal("UserCSS structs don't match.")
	}
}
