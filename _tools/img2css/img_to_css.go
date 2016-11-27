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

// Tool to convert all PNG images in the 'resources/images/' folder
// to CSS as base64 encoded data.

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var out *os.File

const folder = "images/"

func main() {
	// List of images
	images, err := ioutil.ReadDir(folder)
	handleErr(err)

	// Output file
	out, err = os.Create("images.css")
	handleErr(err)
	defer out.Close()

	// Convert all images
	allIcons := ""
	for i, img := range images {
		if !strings.HasSuffix(img.Name(), ".png") {
			continue
		}
		fmt.Print("Processing ", img.Name(), "...")

		iconName := ".gwuimg-" + strings.Split(img.Name(), ".")[0]
		allIcons += iconName
		if i < len(images)-1 {
			allIcons += ", "
		} else {
			allIcons += " "
		}

		out.WriteString(iconName)
		out.WriteString(" {background-image:url(data:image/png;base64,")

		imgin, err := os.Open(folder + img.Name())
		handleErr(err)
		defer imgin.Close()

		buff, err := ioutil.ReadAll(imgin)
		handleErr(err)

		out.WriteString(base64.StdEncoding.EncodeToString(buff))
		out.WriteString(")}\n")

		fmt.Println(" OK")
	}

	// Background formatting for all icons
	out.WriteString("\n")
	out.WriteString(allIcons)
	out.WriteString("{background-position:0px 0px; background-repeat:no-repeat}\n")
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
