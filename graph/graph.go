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
	Name  string
	Links map[string]int
}

var PageMap = map[string]Page{}

func NewPage(name string) Page {
	p := Page{}
	p.Name = name
	p.Links = make(map[string]int)

	util.LogPrint("New page: ", name)

	return p
}

func (p Page) AddLink(path string) {
	_, ok := PageMap[path]
	if !ok {
		PageMap[path] = NewPage(path)
	}

	if _, ok := p.Links[path]; ok {
		p.Links[path]++
	} else {
		p.Links[path] = 1
	}
}
