# Welcome! #

[![GoDoc](https://godoc.org/github.com/icza/gowut/gwu?status.svg)](https://godoc.org/github.com/icza/gowut/gwu) [![Build Status](https://travis-ci.org/icza/gowut.svg?branch=master)](https://travis-ci.org/icza/gowut) [![Go Report Card](https://goreportcard.com/badge/github.com/icza/gowut)](https://goreportcard.com/report/github.com/icza/gowut)

Gowut (Go Web UI Toolkit) is a full-featured, easy to use, platform independent Web UI Toolkit written in pure Go, no platform dependent native code is linked or called.

For documentation please visit the [**Gowut Wiki**](https://github.com/icza/gowut/wiki).

Development takes place in the [`dev` branch](https://github.com/icza/gowut/tree/dev).

## Quick install ##

To quickly install (or update to) the **latest** version, type:

    go get -u github.com/icza/gowut/...

## Quick test ##

To quickly test it and see it in action, run the following example applications.

Let's assume you're in the root of the Gowut project:

    cd $GOPATH/src/github.com/icza/gowut

**1. Showcase of Features.**

This one auto-opens itself in your default browser.

    go run _examples/showcase/showcase.go

The Showcase of Features is also available live: https://gowut-demo.appspot.com/show

**2. A single window example.**

This one auto-opens itself in your default browser.

    go run _examples/simple/simple_demo.go

And this is how it looks:

[![Full App Screenshot](https://github.com/icza/gowut/raw/dev/_images/full_app_example.png)](https://github.com/icza/gowut/wiki/Full-App-Example)

**3. Login window example with session management.**

Change directory so that the demo can read the test cert and key files:

    cd _examples/login
    go run login_demo.go

Open the page `https://localhost:3434/guitest/` in your browser to see it.

## Godoc of Gowut ##

You can read the godoc of Gowut online here:

http://godoc.org/github.com/icza/gowut/gwu

## +1 / Star Gowut! ##
