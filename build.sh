#!/bin/sh
set -ex

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$DIR/cmd/bkp" > /dev/null
rm -f "$(which bkp)"
go install -v
which bkp
popd > /dev/null
