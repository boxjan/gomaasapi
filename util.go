// Copyright 2012-2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package gomaasapi

import (
	"fmt"
	"regexp"
	"strings"
)

var IgnoreTrailUrlSlashRegexPreCompile map[string]*regexp.Regexp
var replaceReg *regexp.Regexp
var DefaultIgnoreTrailSlashUrl = []string{
	"license-key/{osystem}/{distro_series}",
	"scripts/{name}",
	"nodes/{system_id}/blockdevices/{device_id}/partition/{id}",
}

func init() {
	var err error

	// will replace `{.*}` => .*? for url regex string
	replaceReg, _ = regexp.Compile("{.*?}")

	err = ReloadRegex(DefaultIgnoreTrailSlashUrl)
	if err != nil {
		panic(fmt.Errorf("load ignore trail url regex failed with err: %v", err))
	}
}

func ReloadRegex(regexUrl []string) error {

	tmpRegexPreCompile := make(map[string]*regexp.Regexp)
	for _, u := range regexUrl {
		t := replaceReg.ReplaceAllString(u, ".*?")
		tReg, err := regexp.Compile(t + "$")
		// anyone compiled fail, will return error
		if err != nil {
			return fmt.Errorf("compile regix: `%s` failed with err: %+v", t, err)
		}
		tmpRegexPreCompile[u] = tReg
	}

	// all regexp string compile success
	IgnoreTrailUrlSlashRegexPreCompile = tmpRegexPreCompile
	return nil
}

// JoinURLs joins a base URL and a subpath together.
// Regardless of whether baseURL ends in a trailing slash (or even multiple
// trailing slashes), or whether there are any leading slashes at the begining
// of path, the two will always be joined together by a single slash.
func JoinURLs(baseURL, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

// EnsureTrailingSlash appends a slash at the end of the given string unless
// there already is one.
// This is used to create the kind of normalized URLs that Django expects.
// (to avoid Django's redirection when an URL does not ends with a slash.)
func EnsureTrailingSlash(URL string) string {
	for _, rgx := range IgnoreTrailUrlSlashRegexPreCompile {
		if rgx.FindStringIndex(URL) != nil {
			return URL
		}
	}
	if strings.HasSuffix(URL, "/") {
		return URL
	}
	return URL + "/"
}
