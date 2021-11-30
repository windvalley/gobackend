#!/usr/bin/env bash
# run.sh
#

set -e

[[ -z "$1" ]] && {
    echo "Usage: $0 <dev|test|prod>"
    exit 1
}

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
RUN_COMMAND="_output/platforms/$OS/$ARCH/gobackend-apiserver"

cd "$SCRIPT_DIR/../" || exit 1

make gen
make build

echo -e "\n$RUN_COMMAND\n"

RUN_MODE=$1 "$RUN_COMMAND"

exit 0
