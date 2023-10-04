#!/bin/sh
set -o errexit

pkg_with_vuln=()

function scanAuth() {
  cd auth && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
  found=$?
  if [ $found -ne 0 ]; then
   pkg_with_vuln+=("auth")
  fi

  cd ..
}

function scanAwsConfig() {
  cd service/aws/config && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
  found=$?
  if [ $found -ne 0 ]; then
    pkg_with_vuln+=("service/aws/config")
  fi

  cd ../../..
}

function scanAwsSqs() {
  cd service/aws/sqs && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
  found=$?
  if [ $found -ne 0 ]; then
    pkg_with_vuln+=("service/aws/sqs")
  fi

  cd ../../..
}

scanAuth
scanAwsConfig
scanAwsSqs

if [[ -n "${my_array[*]}" ]]; then
  echo "We found vulnerabilities in the following packages ${pkg_with_vuln[@]}"
  echo ""
  echo "Please go to the pkg(s) and update the dependencies, that's one way to solve it."
else
  echo ""
  echo "We not found vulnerabilities"
fi


