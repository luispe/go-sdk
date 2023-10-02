#!/bin/sh
set -o errexit

gofmt -s -w .
go run mvdan.cc/gofumpt@latest -l -w .
