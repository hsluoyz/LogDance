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

import "testing"

func testGetPattern(t *testing.T, path string, res string) {
	t.Helper()
	myRes := GetPattern(path)
	if myRes != res {
		t.Errorf("GetPattern(%s) = %s, supposed to be %s", path, myRes, res)
	}
}

func TestGetPattern(t *testing.T) {
	testGetPattern(t, "/page#tag", "/page")
	testGetPattern(t, "/page/#tag", "/page")

	testGetPattern(t, "/tag/5", "/tag/*")
	testGetPattern(t, "/page/123", "/page/*")

	testGetPattern(t, "/query?id=123", "/query?id=*")

	testGetPattern(t, "/lifestyle-sale/vip/gd2.html?saleType=2&channel=lifestyle&order=s_t_desc&price=0,149&sort=2", "/lifestyle-sale/vip/*.html?saleType=*&channel=*&order=*&price=*&sort=*")
	testGetPattern(t, "/lifestyle-sale/vip/bd1-gd2.html", "/lifestyle-sale/vip/*.html")
}

func testGetDomainName(t *testing.T, url string, res string) {
	t.Helper()
	myRes := GetDomainName(url)
	if myRes != res {
		t.Errorf("GetDomainName(%s) = %s, supposed to be %s", url, myRes, res)
	}
}

func TestGetDomainName(t *testing.T) {
	testGetDomainName(t, "https://www.example.com/", "example.com")
}

func testFormatPath(t *testing.T, url string, domain string, res string) {
	t.Helper()
	myRes := FormatPath(url, domain)
	if myRes != res {
		t.Errorf("getPathFromFullUrl(%s) = %s, supposed to be %s", url, myRes, res)
	}
}

func TestFormatPath(t *testing.T) {
	testFormatPath(t, "//www.example.com/news.aspx", "example.com", "/news.aspx")
	testFormatPath(t, "//travel.example.com/", "example.com", "travel.example.com/")
}

func TestCustomRe(t *testing.T) {
	keyStore["example.com"] = []string{"author", "products"}
	GenerateCustomRe("example.com")

	testGetPattern(t, "/author/alice", "/author/*")
	testGetPattern(t, "/products/abc", "/products/*")
}
