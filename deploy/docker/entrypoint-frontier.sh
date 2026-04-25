#!/bin/bash
# Frontier container entrypoint.
# Renders the config template (FRONTIER_PORT envsubst'd in), waits for the
# shared TLS cert to appear (liaison generates it on first run), then execs
# the frontier binary.

set -eu

CONF_DIR=/opt/liaison/conf
CERTS_DIR=/opt/liaison/certs

: "${FRONTIER_PORT:=30012}"

if [ ! -f "$CONF_DIR/frontier.yaml" ]; then
    export FRONTIER_PORT
    # shellcheck disable=SC2016
    envsubst '${FRONTIER_PORT}' \
        < "$CONF_DIR/frontier.yaml.template" > "$CONF_DIR/frontier.yaml"
    echo "[entrypoint] rendered $CONF_DIR/frontier.yaml (frontier_port=$FRONTIER_PORT)"
fi

for _ in $(seq 1 60); do
    if [ -f "$CERTS_DIR/server.crt" ] && [ -f "$CERTS_DIR/server.key" ]; then
        break
    fi
    echo "[entrypoint] waiting for TLS cert from liaison..."
    sleep 2
done

if [ ! -f "$CERTS_DIR/server.crt" ]; then
    echo "[entrypoint] ERROR: TLS cert never appeared at $CERTS_DIR/server.crt" >&2
    exit 1
fi

exec "$@"
