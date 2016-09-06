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
	"sync/atomic"
)

// ID is the type of the ids of the components.
type ID int

// Converts an ID to a string.
func (id ID) String() string {
	return strconv.Itoa(int(id))
}

// AtoID converts a string to ID.
func AtoID(s string) (ID, error) {
	id, err := strconv.Atoi(s)

	if err != nil {
		return ID(-1), err
	}
	return ID(id), nil
}

// Component id generation and provider

// Last used value for ID
var lastId = new(int64)

// nextCompId returns a unique component id
func nextCompId() ID {
	return ID(atomic.AddInt64(lastId, 1))
}
