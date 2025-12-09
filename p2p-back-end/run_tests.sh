#!/bin/bash
set -e

COVER_FILE="coverage.out"

echo "🔹 Generating list of packages for coverage (exclude entities & servers)..."
coverpkgs=$(go list ./modules/... | grep -vE 'modules/(entities|servers)' | tr '\n' ',' | sed 's/,$//')

echo "🔹 Generating list of packages to test (exclude entities & servers, only with *_test.go)..."
testpkgs=""
for pkg in $(go list ./modules/... | grep -vE 'modules/(entities|servers)'); do
  fs_path=$(echo "$pkg" | sed 's|^p2p-back-end/||') # แปลงเป็น relative path ของ filesystem
  if ls "$fs_path"/*_test.go >/dev/null 2>&1; then
    testpkgs="$testpkgs $pkg"
  fi
done

echo "🔹 Running tests with coverage..."
go test -v -coverpkg=$coverpkgs -coverprofile=$COVER_FILE $testpkgs

echo "🔹 Opening coverage report..."
go tool cover -html=$COVER_FILE
