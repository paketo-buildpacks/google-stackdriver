#!/usr/bin/env bash

set -euo pipefail

if [[ -d ../go-cache ]]; then
  GOPATH=$(realpath ../go-cache)
  export GOPATH
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/google-application-credentials github.com/paketo-buildpacks/google-stackdriver/cmd/google-application-credentials
GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/paketo-buildpacks/google-stackdriver/cmd/main
ln -fs main bin/build
ln -fs main bin/detect
