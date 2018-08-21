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
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/hsluoyz/logdance/graph"
	"github.com/hsluoyz/logdance/pattern"
	"github.com/hsluoyz/logdance/render"
	"github.com/hsluoyz/logdance/target"
)

func printPage(name string, depth int, id uint32, idx int) {
	fmt.Printf("%s[%d-%d] %s\n", strings.Repeat("  ", depth), id, idx, name)
}

func crawl(targetBase string) {
	domain := pattern.GetDomainName(targetBase)
	pattern.GenerateCustomRe(domain)

	graph.AddPage("/")
	printPage("/", 0, 0, 0)
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//fmt.Printf("a[href]: %s\n", e.Attr("href"))
		//fmt.Printf("path: %s\n", e.Request.URL.Path)

		r := e.Request
		href := e.Attr("href")
		after := r.URL.Path
		if !strings.HasSuffix(after, "/") {
			after += "/"
		}

		// Get index of "a[href]".
		var idx int
		if idxAny := r.Ctx.GetAny(fmt.Sprintf("index-%d", r.ID)); idxAny == nil {
			idx = 0
		} else {
			idx = idxAny.(int) + 1
		}
		r.Ctx.Put(fmt.Sprintf("index-%d", r.ID), idx)

		// Check redirection.
		before := r.Ctx.Get("path")
		if before == "" {
			before = "/"
		}
		if idx == 0 && before != after {
			fmt.Printf("(%s != %s)\n", before, after)
			graph.AddRedirectPage(before, after)
		}

		// Get source from previous target.
		source := r.Ctx.Get("pattern")
		if source == "" {
			source = "/"
		}

		// For breakpoint based on ID and index.
		//if r.ID == 8 && idx == 8 {
		//	println("breakpoint here.")
		//}

		target := pattern.StripDomainName(href, domain)

		// Targets like "http://xxx.com", "mailto:xxx@xxx.com", "#tag" will be ignored.
		if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto:") || strings.HasPrefix(target, "#") {
			return
		}

		// Targets like "images/test.jpg/" will be ignored.
		if !pattern.IsHtml(target) {
			return
		}

		// Convert relative URL to site-absolute URL.
		// e.g., "./directions/index.html/" -> "/survivor/directions/index.html/"
		target = pattern.GetAbsolutePath(r.URL.Path, target)

		// Enforce to add the trailing "/" for each path.
		if !strings.HasSuffix(target, "/") {
			target += "/"
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

		if graph.HasPage(tPattern) {
			status = "already done"
		} else {
			printPage(tPattern, r.Depth, r.ID, idx)
		}

		graph.AddLink(sPattern, tPattern)
		// fmt.Printf("New link: [%s] --> [%s]: %s\n", sPattern, tPattern, status)

		if status == "ok" {
			r.Ctx.Put("path", target)
			r.Ctx.Put("pattern", tPattern)
			r.Visit(href)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Printf("OnRequest: %s\n", r.URL.Path)
	})

	c.OnResponse(func(r *colly.Response) {
		//fmt.Printf("OnResponse: %s\n", r.Request.URL.Path)
	})

	c.Visit(targetBase)
}

func main() {
	crawl(target.Url)
	render.GenerateJson()
}
