#!/bin/sh

set -e
cd cmd/bkp
go install
go build
./bkp "$@"
cd - &> /dev/null
