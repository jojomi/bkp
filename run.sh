#!/bin/sh

set -e
go build
./bkp "$@"
