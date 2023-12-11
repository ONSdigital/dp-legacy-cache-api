#!/bin/bash -eux

pushd dp-legacy-cache-api
  make build
  cp build/dp-legacy-cache-api Dockerfile.concourse ../build
popd
