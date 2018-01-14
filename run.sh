#!/bin/sh

./build.sh
./bkp "$@"
cd - &> /dev/null
