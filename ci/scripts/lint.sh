#!/bin/bash -eux

go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
npm install -g @redocly/cli

pushd dp-legacy-cache-api
  make lint
popd
