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
	ucFail = `/*==UserStyle==
@name
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
	domain = `/*==UserStyle==
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

@-moz-document domain(example.com) {
	:root {}
}`
)

func TestValidationPass(t *testing.T) {
	uc := ParseFromString(ucPass)
	pass, err := BasicMetadataValidation(uc)
	if err != nil {
		t.Fatal("Passed validation has err:", err)
	}
	if pass != true {
		t.Fatal("Expected validation to pass.")
	}
}

func TestValidationFail(t *testing.T) {
	uc := ParseFromString(ucFail)
	fail, err := BasicMetadataValidation(uc)
	if err == nil {
		t.Error(err)
	}
	if fail != false {
		t.Fatal("Expected validation to fail.")
	}
}

func TestSingleDomain(t *testing.T) {
	data := ParseFromString(domain)
	pass := Domain{
		Key:   "domain",
		Value: "example.com",
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
