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
	"encoding/json"
	"io/ioutil"
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
	links map[*Page]int
}

var pageMap = map[string]Page {}

func newPage(name string) Page {
	p := Page{}
	p.name = name
	p.links = make(map[*Page]int)
	return p
}

func (p Page) addLink(path string) {
	_, ok := pageMap[path]
	if !ok {
		pageMap[path] = newPage(path)
	}

	page := pageMap[path]
	if _, ok := p.links[&page]; ok {
		p.links[&page] ++
	} else {
		p.links[&page] = 1
	}
}

type Node struct {
	Id    string `json:"id"`
	Group int    `json:"group"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

func newNode(id string, group int) Node {
	n := Node{}
	n.Id = id
	n.Group = group
	return n
}

func newLink(source string, target string, value int) Link {
	l := Link{}
	l.Source = source
	l.Target = target
	l.Value = value
	return l
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func crawl(targetBase string) {
	pageMap["/"] = newPage("/")
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		source := e.Request.URL.Path
		target := e.Attr("href")

		status := "ok"
		if strings.HasPrefix(target, "http") {
			status = "out of scope"
		} else if _, ok := pageMap[target]; ok {
			status = "already done"
		} else {
			pageMap[source].addLink(target)
		}

		fmt.Printf("New link: [%s] --> [%s]: %s\n", source, target, status)

		if status == "ok" {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("New page: ", r.URL.Path)
	})

	c.Visit(targetBase)
}

func generateJson() {
	g := Graph{}
	g.Nodes = make([]Node, 0)
	g.Links = make([]Link, 0)

	for _, page := range pageMap {
		g.Nodes = append(g.Nodes, newNode(page.name, 0))

		for target, weight := range page.links {
			g.Links = append(g.Links, newLink(page.name, target.name, weight))
		}
	}

	if data, err := json.MarshalIndent(g, "", "  "); err == nil {
		fmt.Printf("%s\n", data)

		err := ioutil.WriteFile("webgraph.json", data, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	targetBase := "http://127.0.0.1:5000/"

	// loadLogFile("log/raith.log", apacheLogHandler)

	crawl(targetBase)
	generateJson()
}
