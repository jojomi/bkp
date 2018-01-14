#!/bin/sh

set -e

bkp=../cmd/bkp/bkp
target_dir=/tmp/bkp-backup/


echo "Building..."
rm -f "$bkp"
../build.sh
echo "Built."


echo "Clearing state..."
rm -rf "$target_dir"
mkdir -p "$target_dir"
echo "Clear."

echo "Starting backup..."
$bkp --config-dirs config --jobs testjob


# rm -rf "$target_dir"