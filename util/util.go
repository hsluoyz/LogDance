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

package util

import (
	"log"
	"os"
)

// EnableLog controls whether to print log to console.
var EnableLog = true

func init() {
	//create your file with desired read/write permissions
	f, err := os.OpenFile("page.log", os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	// defer f.Close()

	//set output of logs to f
	log.SetOutput(f)
}

// LogPrint prints the log.
func LogPrint(v ...interface{}) {
	if EnableLog {
		log.Print(v...)
	}
}

// LogPrintf prints the log with the format.
func LogPrintf(format string, v ...interface{}) {
	if EnableLog {
		log.Printf(format, v...)
	}
}
