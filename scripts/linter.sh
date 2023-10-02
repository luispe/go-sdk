#!/bin/sh
set -o errexit

go run github.com/mgechev/revive@latest -config metalint.toml -formatter friendly ./...
