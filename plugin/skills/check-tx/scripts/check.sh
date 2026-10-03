#!/usr/bin/env bash
set -euo pipefail

shapes="${1:?usage: check.sh <shapes.json> [isolation] [txns]}"
iso="${2:-read-committed}"
txns="${3:-300}"
image="${ELLE_GO_IMAGE:-ghcr.io/syed-imam/elle-go:latest}"

[ -r "$shapes" ] || { echo "cannot read shapes file: $shapes" >&2; exit 2; }
docker info >/dev/null 2>&1 || { echo "Docker is not running" >&2; exit 2; }

name="elle-check-$$"
docker run -d --rm --name "$name" "$image" >/dev/null
trap 'docker stop "$name" >/dev/null 2>&1 || true' EXIT

for _ in $(seq 1 120); do
  docker exec "$name" pg_isready -h 127.0.0.1 -U postgres >/dev/null 2>&1 && break
  sleep 0.5
done

set +e
docker exec -i "$name" elle-go pg \
  -dsn "postgres://postgres:elle@127.0.0.1:5432/postgres?sslmode=disable" \
  -isolation "$iso" -shapes /dev/stdin -txns "$txns" < "$shapes"
code=$?
set -e
exit "$code"
