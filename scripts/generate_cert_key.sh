#!/usr/bin/env bash
# generate_cert_key.sh
#
# For generating a development cert and key

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

cd "$SCRIPT_DIR" || exit 1

CERT_DIR=$SCRIPT_DIR/../configs/cert
DOMAIN="localhost"
GOROOT="$(go env | grep GOROOT | awk -F\" '{print $2}')"

mkdir -p "$CERT_DIR"

(
    cd "$CERT_DIR" || exit 1

    go run "$GOROOT"/src/crypto/tls/generate_cert.go --host="$DOMAIN"
)

exit 0
