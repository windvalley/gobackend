#!/usr/bin/env bash
# run_dev.sh
#

set -e

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OS=$(go env GOOS)
ARCH=$(go env GOARCH)

cd "$SCRIPT_DIR/../" || exit 1

make build

_output/platforms/"$OS"/"$ARCH"/gobackend-apiserver -c configs/dev.gobackend-apiserver.yaml

exit 0
