package main

import "testing"

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
