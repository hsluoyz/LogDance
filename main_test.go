package main

import "testing"

func testGetPattern(t *testing.T, path string, res string) {
	t.Helper()
	myRes := getPattern(path)
	t.Logf("getPattern(%s) = %s", path, myRes)

	if myRes != res {
		t.Errorf("getPattern(%s) = %s, supposed to be %s", path, myRes, res)
	}
}

func TestGetPattern(t *testing.T) {
	testGetPattern(t, "/page#tag", "/page")
	testGetPattern(t, "/page/#tag", "/page")
	testGetPattern(t, "/tag/5", "/tag/*")
	testGetPattern(t, "/page/123", "/page/*")
	testGetPattern(t, "/author/alice", "/author/*")
	testGetPattern(t, "/products/abc", "/products/*")
	testGetPattern(t, "/query?id=123", "/query?id=*")
	testGetPattern(t, "", "")
	testGetPattern(t, "/lifestyle-sale/vip/gd2.html?saleType=2&channel=lifestyle&order=s_t_desc&price=0,149&sort=2", "/lifestyle-sale/vip/gd2.html?saleType=*&channel=*&order=*&price=*&sort=*")
}

func testGetDomainName(t *testing.T, url string, res string) {
	t.Helper()
	myRes := getDomainName(url)
	t.Logf("getDomainName(%s) = %s", url, myRes)

	if myRes != res {
		t.Errorf("getDomainName(%s) = %s, supposed to be %s", url, myRes, res)
	}
}

func TestGetDomainName(t *testing.T) {
	testGetDomainName(t, "https://www.example.com/", "example.com")
}

func testFormatPath(t *testing.T, url string, domain string, res string) {
	t.Helper()
	myRes := formatPath(url, domain)
	t.Logf("getPathFromFullUrl(%s) = %s", url, myRes)

	if myRes != res {
		t.Errorf("getPathFromFullUrl(%s) = %s, supposed to be %s", url, myRes, res)
	}
}

func TestFormatPath(t *testing.T) {
	testFormatPath(t, "//www.example.com/news.aspx", "example.com", "/news.aspx")
	testFormatPath(t, "//travel.example.com/", "example.com", "travel.example.com/")
}
