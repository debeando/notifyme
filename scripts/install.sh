#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

FILE=
TAG=$(curl --silent "https://api.github.com/repos/swapbyt3s/notifyme/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -f /usr/local/bin/notifyme ]; then
  rm -f /usr/bin/notifyme
  rm -f /usr/local/bin/notifyme
fi

if [[ "${OSTYPE}" == "darwin"* ]]; then
  FILE="notifyme-darwin_amd64.tar.gz"
elif [[ "${OSTYPE}" == "linux"* ]]; then
  FILE="notifyme-linux_amd64.tar.gz"
fi

if [ ! -z "${FILE}" ]; then
  wget -qO- https://github.com/swapbyt3s/notifyme/releases/download/${TAG}/${FILE} | tar xz -C /usr/local/bin/
fi

if [ -f /usr/local/bin/notifyme ]; then
  if [[ "${OSTYPE}" == "linux"* ]]; then
    ln -s /usr/local/bin/notifyme /usr/bin/notifyme
  fi
fi
