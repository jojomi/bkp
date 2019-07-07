#!/bin/sh
set -ex

echo "Building..."
./build.sh
echo "Running..."
sudo bkp "$@"
