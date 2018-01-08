#!/bin/sh

set -e
cd ui
go build
./bkp "$@"
cd -
