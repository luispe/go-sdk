#!/bin/bash
set -eo pipefail

execute_script_in_directory() {
  cd $1
  go run honnef.co/go/tools/cmd/staticcheck@latest -f stylish ./...

  SCRIPT_DIR=$(dirname $0)
  cd $SCRIPT_DIR/..
}

execute_script_in_directory "$1"

echo ""
echo "Congratulations, we didn't find any improvements for your code"
exit 0
