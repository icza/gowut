// +build appengine

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

// Implementation of the GUI server Start on Google App Engine.

package gwu

import (
	"log"
	"net/http"
)

func (s *serverImpl) Start(openWins ...string) error {
	http.HandleFunc(s.appPath, func(w http.ResponseWriter, r *http.Request) {
		s.serveHTTP(w, r)
	})

	http.HandleFunc(s.appPath+pathStatic, func(w http.ResponseWriter, r *http.Request) {
		s.serveStatic(w, r)
	})

	log.Println("GAE - Starting GUI server on path:", s.appPath)
	if s.logger != nil {
		s.logger.Println("GAE - Starting GUI server on path:", s.appPath)
	}

	go s.sessCleaner()

	return nil
}
