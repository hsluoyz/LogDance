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
	"regexp"
	"strings"
)

func GetPattern(path string) string {
	// "/page#tag" -> "/page"
	// "/page/#tag" -> "/page"
	re, _ := regexp.Compile("/?#.*")
	path = re.ReplaceAllString(path, "")

	//// "/author/alice" -> "/author/*"
	//re, _ = regexp.Compile("(author/)[^/]*(.*)")
	//path = re.ReplaceAllString(path, "$1*$2")
	//
	//// "/products/abc" -> "/products/*"
	//re, _ := regexp.Compile("(products/)[^/]*(.*)")
	//path = re.ReplaceAllString(path, "$1*$2")

	// "/query?id=123" -> "/query?id=*"
	re, _ = regexp.Compile("=[^&=]*")
	path = re.ReplaceAllString(path, "=*")

	// "/page5" -> "/page*"
	re, _ = regexp.Compile("[0-9]+")
	path = re.ReplaceAllString(path, "*")

	// "/2018/01/02/the-blog-title" -> "/xxxx/xx/xx/xx"
	re, _ = regexp.Compile("(.*)\\d{4}/\\d{2}/\\d{2}/(.*)")
	path = re.ReplaceAllString(path, "$1xxxx/xx/xx/xx")

	return path
}

func GetDomainName(url string) string {
	re, _ := regexp.Compile("[^.]*\\.(com|net|org)")
	url = re.FindString(url)
	return url
}

func GetSubDomain(url string, domain string) string {
	if i := strings.Index(url, domain); i != -1 {
		url = url[:i]
		re, _ := regexp.Compile("/[^./]*\\.")
		url = re.FindString(url)
		if url == "" {
			return ""
		} else {
			return url[1:len(url) - 1]
		}
	} else {
		return ""
	}
}

func FormatPath(url string, domain string) string {
	if i := strings.Index(url, domain); i != -1 {
		subDomain := GetSubDomain(url, domain)
		if subDomain == "" || subDomain == "www" {
			return url[i + len(domain):]
		} else {
			return strings.TrimLeft(url, "/")
		}
	} else {
		return url
	}
}
