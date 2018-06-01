#!/usr/bin/env bash

set -e

cd ../

go build -ldflags "-s -w" -o nm main.go
