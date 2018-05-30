#!/usr/bin/env bash

set -e

cd ../

go build -ldflags "-s -w" -o notifyme main.go
