// +build !appengine

// Copyright (C) 2013 Andras Belicza. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Contains the main function of the Gowut "Showcase of Features" demo.
// Separated because main() can't be defined on AppEngine.

package main

import (
	"flag"

	"github.com/icza/gowut/_examples/showcase/showcasecore"
)

var (
	addr     = flag.String("addr", "", "address to start the server on")
	appName  = flag.String("appName", "showcase", "Gowut app name")
	autoOpen = flag.Bool("autoOpen", true, "auto-open the demo in default browser")
)

func main() {
	flag.Parse()

	showcasecore.StartServer(*appName, *addr, *autoOpen)
}
