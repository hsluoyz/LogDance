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

package main

import (
	"os"
	"bufio"
	"strings"
	"github.com/Songmu/axslogparser"
	"github.com/gocolly/colly"
	"fmt"
)

func loadLogFile(filePath string, handler func(string)) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		handler(line)
	}
	return scanner.Err()
}

func apacheLogHandler(line string) {
	if line == "" {
		return
	}

	l, err := axslogparser.Parse(line)
	if err != nil {
		panic(err)
	}
	printApacheLog(l)
}

func printApacheLog(l *axslogparser.Log) {
	println(l.Host, l.Time.String(), l.Request, l.Status, l.Referer, l.UserAgent, l.RequestURI, l.Method)
}

type Page struct {
	name string
	links []*Page
}

var pageMap = map[string]Page {}

func newPage(name string) Page {
	page := Page{}
	page.name = name
	return page
}

func crawl(targetBase string) {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		source := e.Request.URL.Path
		target := e.Attr("href")
		if !strings.HasPrefix(target, "http") {
			fmt.Printf("New link: [%s] --> [%s]: ok\n", source, target)
			e.Request.Visit(e.Attr("href"))
		} else {
			fmt.Printf("New link: [%s] --> [%s]: deny\n", source, target)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("New page: ", r.URL.Path)

		_, ok := pageMap[r.URL.Path]
		if !ok {
			pageMap[r.URL.Path] = newPage(r.URL.Path)
		}
	})

	c.Visit(targetBase)
}

func main() {
	targetBase := "http://127.0.0.1:5000/"

	// loadLogFile("log/raith.log", apacheLogHandler)

	crawl(targetBase)
	print(pageMap)
}
