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

package render

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hsluoyz/logdance/graph"
)

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

func GenerateJson() {
	g := Graph{}
	g.Nodes = make([]Node, 0)
	g.Links = make([]Link, 0)

	for i, page := range graph.PageList {
		page.Id = i
	}

	for _, page := range graph.PageList {
		if page.Name == "/" {
			g.Nodes = append(g.Nodes, newNode(page.Id, page.Name, "home"))
		} else {
			g.Nodes = append(g.Nodes, newNode(page.Id, page.Name, "page"))
		}

		for target := range page.Links {
			g.Links = append(g.Links, newLink(page.Id, graph.PageMap[target].Id))
		}
	}

	if data, err := json.MarshalIndent(g, "", "  "); err == nil {
		// fmt.Printf("%s\n", data)

		err := ioutil.WriteFile("webgraph.json", data, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
