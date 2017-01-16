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

// Implementation of the GUI server Start in standalone apps (non-GAE).

package gwu

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func (s *serverImpl) Start(openWins ...string) error {
	http.HandleFunc(s.appPath, func(w http.ResponseWriter, r *http.Request) {
		s.serveHTTP(w, r)
	})

	http.HandleFunc(s.appPath+pathStatic, func(w http.ResponseWriter, r *http.Request) {
		s.serveStatic(w, r)
	})

	appURL := s.AppURL()
	log.Println("Starting GUI server on:", appURL)
	if s.logger != nil {
		s.logger.Println("Starting GUI server on:", appURL)
	}

	for _, winName := range openWins {
		if err := open(appURL + winName); err != nil {
			s.logger.Printf("Opening window '%s' err: %v\n", appURL+winName, err)
		}
	}

	go s.sessCleaner()

	var err error
	if s.secure {
		err = http.ListenAndServeTLS(s.addr, s.certFile, s.keyFile, nil)
	} else {
		err = http.ListenAndServe(s.addr, nil)
	}

	if err != nil {
		return err
	}
	return nil
}
