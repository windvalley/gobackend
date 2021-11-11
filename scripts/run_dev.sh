#!/usr/bin/env bash
# run_dev.sh
#

set -e

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
RUN_COMMAND="_output/platforms/$OS/$ARCH/gobackend-apiserver -c configs/dev.gobackend-apiserver.yaml"

cd "$SCRIPT_DIR/../" || exit 1

make gen
make build

echo -e "\n$RUN_COMMAND\n"

eval "$RUN_COMMAND"

exit 0
