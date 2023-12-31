#!/bin/bash
set -eo pipefail

execute_script_in_directory() {
  cd $1
  go test -race $(go list ./... | grep -v mocks) && cd ../../..

  SCRIPT_DIR=$(dirname $0)
  cd $SCRIPT_DIR/..
}

execute_script_in_directory "$1"
