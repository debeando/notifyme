#!/usr/bin/env bash

cd ../

export TAG=$1

if [[ -z $TAG ]]
then
  echo "No tag defined. E.g.: '0.1.0'"
  exit 1
fi

git push --delete origin "v${TAG}"
git tag "v${TAG}"
git push --tags --force

curl --data "{\"tag_name\": \"v${TAG}\",\"target_commitish\": \"master\",\"name\": \"Pre Release v${TAG}\",\"body\": \"\",\"draft\": true,\"prerelease\": true}" https://api.github.com/repos/swapbyt3s/NotifyMe/releases?access_token=${GITHUB_TOKEN}

ID=$(curl -sH "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/swapbyt3s/NotifyMe/releases | jq -r '.[0].id')

rm -rf pkg/*

declare -a OS=("linux" "darwin")

for os in "${OS[@]}"
do
  mkdir -p pkg/${os}_amd64/
  GOOS=${os} GOARCH=amd64 go build -ldflags "-s -w" -o pkg/${os}_amd64/notifyme main.go
  tar -czvf pkg/${os}_amd64/notifyme-${TAG}-${os}_amd64.tar.gz pkg/${os}_amd64/notifyme

  curl -# \
       -XPOST \
       -H "Authorization:token ${GITHUB_TOKEN}" \
       -H "Content-Type:application/octet-stream" \
       --data-binary @pkg/${os}_amd64/notifyme-${TAG}-${os}_amd64.tar.gz \
       https://uploads.github.com/repos/swapbyt3s/NotifyMe/releases/${ID}/assets?name=notifyme-${TAG}-${os}_amd64.tar.gz
done
