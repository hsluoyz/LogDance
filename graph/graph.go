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

package graph

import "github.com/hsluoyz/logdance/util"

type Page struct {
	Id    int
	Name  string
	Links map[string]int
}

var PageMap = map[string]*Page{}

func newPage(name string) *Page {
	p := Page{}
	p.Name = name
	p.Links = make(map[string]int)

	util.LogPrint("New page: ", name)

	return &p
}

func (p Page) addLink(path string) {
	_, ok := PageMap[path]
	if !ok {
		PageMap[path] = newPage(path)
	}

	if _, ok := p.Links[path]; ok {
		p.Links[path]++
	} else {
		p.Links[path] = 1
	}
}

func AddPage(name string) {
	PageMap[name] = newPage(name)
}

// For redirect: "/" -> "/home.html/",
// before = "/"
// after = "/home.html/"
func AddRedirectPage(before string, after string) {
	page, ok := PageMap[before]
	if !ok {
		panic("\"before\" of a redirection does not exist")
	}

	// Maybe:
	// before = "/home.html/"
	// after = "/"
	// So we use the shorter one from before and after as the final name.
	if len(after) < len(before) {
		page.Name = after
	}

	PageMap[after] = page
}

func HasPage(name string) bool {
	_, ok := PageMap[name]
	return ok
}

func AddLink(sPage string, tPage string) {
	PageMap[sPage].addLink(tPage)
}
