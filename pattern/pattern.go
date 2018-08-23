// Copyright (c) Microsoft Corporation. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pattern

import (
	"net/url"
	"regexp"
	"strings"
)

var keyStore map[string][]string
var customRe []*regexp.Regexp

func init() {
	keyStore = make(map[string][]string)

	keyStore["quotes.toscrape.com"] = []string{"author|tag"}
	keyStore["gaohaoyang.github.io"] = []string{"\\d{4}/\\d{2}/\\d{2}"} // "/2018/01/02/the-blog-title" -> "/*/*/*/*"
	keyStore["yohobuy.com"] = []string{"shop|tags"}
	keyStore["www.ruanyifeng.com"] = []string{"blog|survivor|road"}
	keyStore["books.toscrape.com"] = []string{"books", "catalogue"}
}

func GetPattern(path string) string {
	// "/page#tag" -> "/page"
	// "/page/#tag" -> "/page"
	re, _ := regexp.Compile("/?#.*")
	path = re.ReplaceAllString(path, "/")

	//// "/author/alice" -> "/author/*"
	//re, _ = regexp.Compile("(author/)[^/]*(.*)")
	//path = re.ReplaceAllString(path, "$1*$2")
	//
	//// "/products/abc" -> "/products/*"
	//re, _ := regexp.Compile("(products/)[^/]*(.*)")
	//path = re.ReplaceAllString(path, "$1*$2")

	for _, re := range customRe {
		path = re.ReplaceAllString(path, "$1/*$2")
	}

	// "/query?id=123" -> "/query?id=*"
	re, _ = regexp.Compile("=[^&=]*")
	path = re.ReplaceAllString(path, "=*")

	// "/page5" -> "/page*"
	re, _ = regexp.Compile("[0-9]+")
	path = re.ReplaceAllString(path, "*")

	// "/products/abc.html" -> "/products/*.html"
	if strings.Contains(path, "*") {
		re, _ = regexp.Compile("(.*/)[^./]*(.html.*)")
		path = re.ReplaceAllString(path, "$1*$2")
	}

	return path
}

func GetFullDomainName(url string) string {
	i := strings.Index(url, "//")
	if i == -1 {
		panic("GetDomainName() error: no \"//\" in url")
	}
	i += 2

	j := len(url)
	if url[len(url) - 1] == '/' {
		j --
	}

	return url[i:j]
}

func GetDomainName(url string) string {
	full := GetFullDomainName(url)

	i := strings.LastIndex(full, ".")
	if i == -1 {
		panic("GetDomainName() error: no \".\" in url")
	}

	j := strings.LastIndex(full[:i], ".")
	if j == -1 {
		return full
	} else {
		return full[j+1:]
	}
}

func GetSubDomain(url string, domain string) string {
	if i := strings.Index(url, domain); i != -1 {
		url = url[:i]
		re, _ := regexp.Compile("/[^./]*\\.")
		url = re.FindString(url)
		if url == "" {
			return ""
		} else {
			return url[1 : len(url)-1]
		}
	} else {
		return ""
	}
}

func StripDomainName(url string, domain string) string {
	if i := strings.Index(url, domain); i != -1 {
		subDomain := GetSubDomain(url, domain)
		if subDomain == "" || subDomain == "www" {
			return url[i+len(domain):]
		} else {
			return strings.TrimLeft(url, "/")
		}
	} else {
		return url
	}
}

// domain is like "/author/alice"
// regex is like "(author)/[^/]+(.*)"
// replaced with "$1/*$2"
func GenerateCustomRe(fullDomain string) {
	if keys, ok := keyStore[fullDomain]; ok {
		for _, key := range keys {
			expr := "(" + key + ")/[^/]+(.*)"
			re, _ := regexp.Compile(expr)
			customRe = append(customRe, re)
		}
	} else {
		customRe = nil
	}
}

func GetAbsolutePath(base string, path string) string {
	p, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	b, err := url.Parse(base)
	if err != nil {
		panic(err)
	}

	res := b.ResolveReference(p)
	if res.RawQuery == "" {
		return res.Path
	} else {
		return res.Path + "?" + res.RawQuery
	}
}

func IsHtml(path string) bool {
	i := strings.LastIndex(path, ".")
	if i == -1 {
		return true
	} else {
		ext := path[i+1:]
		if strings.HasPrefix(ext, "htm") {
			return true
		} else {
			return false
		}
	}
}
