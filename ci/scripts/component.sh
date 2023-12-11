#!/bin/bash -eux

pushd dp-legacy-cache-api
  make test-component
popd
