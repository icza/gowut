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

// ID type definition, and unique ID generation.

package gwu

import (
	"strconv"
)

// The type of the ids of the components.
type ID int

// Converts an ID to a string.
func (id ID) String() string {
	return strconv.Itoa(int(id))
}

// Converts a string to ID
func AtoID(s string) (ID, error) {
	id, err := strconv.Atoi(s)

	if err != nil {
		return ID(-1), err
	}
	return ID(id), nil
}

// Component id generation and provider

// A channel used to generate unique ids
var idChan chan ID = make(chan ID)

// init stats a new go routine to generate unique ids
func init() {
	go func() {
		for i := 0; ; i++ {
			idChan <- ID(i)
		}
	}()
}

// nextCompId returns a unique component id
func nextCompId() ID {
	return <-idChan
}
