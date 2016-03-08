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
	"fmt"
	"github.com/icza/gowut/examples/showcase/showcasecore"
	"log"
	"os"
)

func main() {
	// Allow app control from command line (in co-operation with the starter script):
	log.Println("Type 'r' to restart, 'e' to exit.")
	go func() {
		var cmd string
		for {
			fmt.Scanf("%s", &cmd)
			switch cmd {
			case "r": // restart
				os.Exit(1)
			case "e": // exit
				os.Exit(0)
			}
		}
	}()

	showcasecore.StartServer("showcase")
}
