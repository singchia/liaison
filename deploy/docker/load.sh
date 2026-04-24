#!/bin/bash
# Load the bundled Liaison Docker images from images/*.tar into the local daemon.
# Idempotent — re-running is safe.
set -eu

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
IMAGES_DIR="$SCRIPT_DIR/images"

if [ ! -d "$IMAGES_DIR" ]; then
    echo "❌ $IMAGES_DIR not found. Are you running this from the extracted package root?" >&2
    exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
    echo "❌ docker not found on PATH. Install Docker 20.10+ first." >&2
    exit 1
fi

found=0
for tar in "$IMAGES_DIR"/*.tar; do
    [ -f "$tar" ] || continue
    found=1
    echo "📥 Loading $(basename "$tar")..."
    docker load -i "$tar"
done

if [ "$found" -eq 0 ]; then
    echo "❌ No *.tar images found under $IMAGES_DIR." >&2
    exit 1
fi

echo
echo "✅ Images loaded. Next steps:"
echo "   1. cp .env.example .env   # edit LIAISON_PUBLIC_HOST etc."
echo "   2. docker compose up -d"
echo "   3. docker compose logs liaison | grep -A5 'first-run credentials'"
