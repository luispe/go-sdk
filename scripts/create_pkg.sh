#!/bin/sh
set -o errexit

function createPkg() {
    mkdir -p "$1" \
    && cd "$1" \
    && go mod init "github.com/pomelo-la/go-toolkit/${1}" \
    && go mod tidy
}

PKG_NAME=$(gum input --placeholder "package name")

gum confirm "We create a go package with name: ${PKG_NAME}
Save changes?" && createPkg ${PKG_NAME}
