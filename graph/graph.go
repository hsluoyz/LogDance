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
	Id      int         `json:"nodes"`
	Name    string      `json:"nodes"`
	Aliases []string    `json:"nodes"`
	Links   map[int]int `json:"nodes"`
}

var PageList = []*Page{}
var PageMap = map[string]*Page{}

func newPage(id int, name string) *Page {
	p := Page{}
	p.Id = id
	p.Name = name
	p.Links = make(map[int]int)

	util.LogPrint("New page: ", name)

	return &p
}

func (p *Page) addAlias(name string) {
	p.Aliases = append(p.Aliases, name)
}

func (p *Page) addLink(path string) {
	target, ok := PageMap[path]
	if !ok {
		target = newPage(len(PageList), path)
		PageList = append(PageList, target)
		PageMap[path] = target
	}

	if _, ok := p.Links[target.Id]; ok {
		p.Links[target.Id]++
	} else {
		p.Links[target.Id] = 1
	}
}

func AddPage(name string) {
	newPage := newPage(len(PageList), name)
	PageList = append(PageList, newPage)
	PageMap[name] = newPage
}

// For redirect: "/" -> "/home.html/",
// before = "/"
// after = "/home.html/"
func AddRedirectPage(before string, after string) {
	page, ok := PageMap[before]
	if !ok {
		panic("\"before\" of a redirection does not exist")
	}

	afterPage, ok := PageMap[after]
	if !ok {
		// Maybe:
		// before = "/home.html/"
		// after = "/"
		// So we use the shorter one from before and after as the final name.
		if len(after) < len(before) {
			page.Name = after
		}

		page.addAlias(after)
		PageMap[after] = page
	} else {
		i := page.Id
		// Delete the before page (i-th) because the after page already exists.
		PageList = append(PageList[: i], PageList[i + 1 :]...)
		afterPage.addAlias(before)
		PageMap[before] = page
	}
}

func HasPage(name string) bool {
	_, ok := PageMap[name]
	return ok
}

func AddLink(sPage string, tPage string) {
	PageMap[sPage].addLink(tPage)
}
