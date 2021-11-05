#!/usr/bin/env bash
# change_project_name.sh
#
# Change current project name 'go-web-backend' to your own project name.
#
# e.g.:
#   ./change_project_name.sh your-project-name

set -e

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

cd "$SCRIPT_DIR" || exit 1

[[ -z "$1" ]] && {
    echo Usage: "$0" new_project_name
    exit 1
}

CURRENT_PROJECT_NAME=$(awk 'NR==1{print $2}' ../go.mod)
NEW_PROJECT_NAME=$1
SED="sed"

[[ $(uname) = "Darwin" ]] && SED=gsed

# shellcheck disable=SC2038
find ../ -type f -name "*.go" -o -name "go.mod" |
    xargs $SED -i "s#${CURRENT_PROJECT_NAME}#${NEW_PROJECT_NAME}#"

exit 0
