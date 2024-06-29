#!/bin/env bash

set -eu
set -o pipefail

distListJson=$(go tool dist list -json)

if [ -d "bin" ]; then
    rm -rf bin
fi
mkdir bin

for row in $(echo "${distListJson}" | jq -r '.[] | @base64'); do
    dist=$(echo "$row" | base64 --decode | jq -r '.GOOS')
    arch=$(echo "$row" | base64 --decode | jq -r '.GOARCH')

    if [[ $dist == "linux" && $arch == a* || $dist == "darwin" && $arch == a* ]]; then
        GOOS=$dist GOARCH=$arch go build --ldflags '-extldflags "-static"' -o bin/loadbalancer cmd/lb.go
        cd bin
        tar -czvf "loadbalancer-$dist-$arch.tar.gz" loadbalancer
        rm loadbalancer
        cd ..
    fi
done
