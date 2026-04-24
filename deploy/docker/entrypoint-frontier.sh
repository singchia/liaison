#!/bin/bash
# Frontier container entrypoint.
# Renders the config template from env vars, then execs the frontier binary.
# Certs are provided by the shared volume populated by the liaison container
# (liaison generates the cert pair on first run). Frontier blocks until the
# cert exists so the two boot orders are tolerated.

set -eu

CONF_DIR=/opt/liaison/conf
CERTS_DIR=/opt/liaison/certs

if [ ! -f "$CONF_DIR/frontier.yaml" ]; then
    cp "$CONF_DIR/frontier.yaml.template" "$CONF_DIR/frontier.yaml"
    echo "[entrypoint] seeded $CONF_DIR/frontier.yaml from template"
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
