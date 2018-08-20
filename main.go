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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/hsluoyz/logdance/pattern"
	"github.com/hsluoyz/logdance/util"
)

type Page struct {
	name  string
	links map[string]int
}

var pageMap = map[string]Page{}

func newPage(name string) Page {
	p := Page{}
	p.name = name
	p.links = make(map[string]int)

	util.LogPrint("New page: ", name)

	return p
}

func (p Page) addLink(path string) {
	_, ok := pageMap[path]
	if !ok {
		pageMap[path] = newPage(path)
	}

	if _, ok := p.links[path]; ok {
		p.links[path]++
	} else {
		p.links[path] = 1
	}
}

type Node struct {
	Id    string `json:"id"`
	Group int    `json:"group"`
	Size  int    `json:"size"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

func newNode(id string, group int, size int) Node {
	n := Node{}
	n.Id = id
	n.Group = group
	n.Size = size
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

func printPage(name string, depth int, id uint32) {
	fmt.Printf("%s[%d] %s\n", strings.Repeat(" ", depth), id, name)
}

func crawl(targetBase string) {
	domain := pattern.GetDomainName(targetBase)
	pattern.GenerateCustomRe(domain)

	pageMap["/"] = newPage("/")
	printPage("/", 0, 0)
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		source := e.Request.URL.Path
		target := e.Attr("href")

		if e.Request.ID == 1 {
			source = "/"
		} else {
			if source != e.Request.Ctx.Get("path") {
				// println(source + " != " + e.Request.Ctx.Get("path"))
			}
			if origin := e.Request.Ctx.Get("pattern"); origin != "" {
				source = origin
			} else {
				if e.Request.URL.RawQuery != "" {
					source += "?" + e.Request.URL.RawQuery
				}
				if source == "" {
					source = "/"
				}
				if source != "/" {
					source = strings.TrimSuffix(source, "/")
				}
			}
		}

		target = pattern.FormatPath(target, domain)

		if target != "/" {
			target = strings.TrimSuffix(target, "/")
		}

		// Targets like "http://xxx.com", "mailto:xxx@xxx.com", "#tag" will be ignored.
		if target == "" || strings.HasPrefix(target, "..") || strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") || strings.HasPrefix(target, "#") {
			return
		}

		status := "ok"
		sPattern := pattern.GetPattern(source)
		tPattern := pattern.GetPattern(target)
		if sPattern == tPattern {
			return
		}
		// Do not handle the main page again by recognizing "index.htm".
		if strings.HasPrefix(tPattern, "/index.htm") {
			return
		}

		if _, ok := pageMap[tPattern]; ok {
			status = "already done"
		}
		pageMap[sPattern].addLink(tPattern)
		printPage(tPattern, e.Request.Depth, e.Request.ID)

		// fmt.Printf("New link: [%s] --> [%s]: %s\n", sPattern, tPattern, status)

		if status == "ok" {
			e.Request.Ctx.Put("path", target)
			e.Request.Ctx.Put("pattern", tPattern)
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.OnRequest(func(r *colly.Request) {
	})

	c.Visit(targetBase)
}

func generateJson() {
	g := Graph{}
	g.Nodes = make([]Node, 0)
	g.Links = make([]Link, 0)

	for _, page := range pageMap {
		if page.name == "/" {
			g.Nodes = append(g.Nodes, newNode(page.name, 0, 20))
		} else {
			g.Nodes = append(g.Nodes, newNode(page.name, 0, 10))
		}

		for target := range page.links {
			g.Links = append(g.Links, newLink(page.name, target, 1))
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
	targetBase := "http://www.ruanyifeng.com"

	crawl(targetBase)
	generateJson()
}
