#!/usr/bin/env bash

set -e

go build -ldflags "-s -w" -o ../notifyme ../main.go
