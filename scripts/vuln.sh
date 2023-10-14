#!/bin/bash
set -eo pipefail

execute_script_in_directory() {
  cd $1
  go run golang.org/x/vuln/cmd/govulncheck@latest ./...
  FOUND=$?
    if [ $FOUND -ne 0 ]; then
       echo "We found vulnerabilities in the following packages ${1}"
       echo ""
       echo "Please go to the pkg(s) and update the dependencies, that's one way to solve it."
       exit 1
    fi
  SCRIPT_DIR=$(dirname $0)
  cd $SCRIPT_DIR/..
}

execute_script_in_directory "$1"

echo ""
echo "Congratulations, we didn't find any vulnerabilities for your code"
exit 0
