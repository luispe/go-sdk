# Getting started for Developers

## Contributing

### Before You Start
The project is written in Golang. If you do not have a 
good grounding in Go, try out [the tutorial](https://tour.golang.org/).

!!! info

    One of the main concepts for working with go-toolkit is that each
    module/package works separately, this makes administration easier
    and for the client it is a better experience since it only
    downloads the transitive dependencies of a single module.

## Pre-requisites

- Go installed (see [install guide](https://go.dev/dl))

The directory structure is opinionated and idiomatic, 
we maintain a flat structure where we try to make 
each package self-descriptive from its name.

Each package/module has its own go.mod and go.sum 
getting the advantages of working in the same repository 
but with separate dependency management for each package.

## Format code

    make gofmt

## Scan vulnerabilities

    make vuln PKG_NAME=<PACKAGE_NAME>

e.g:

    make vuln PKG_NAME=httprouter

## Running test

    make test PKG_NAME=<PACKAGE_NAME>

e.g:

    make test PKG_NAME=httprouter

## Running lint

    make lint

## Running static check

    make static_checks PKG_NAME=<PACKAGE_NAME>

e.g:

    make test PKG_NAME=httprouter

## Documentation Changes

Modify contents in docs/ directory.

Preview changes in your browser by visiting http://localhost:8000 after running:

    make serve-docs
