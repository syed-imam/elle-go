#!/usr/bin/env bash
set -euo pipefail

shapes="${1:?usage: check.sh <shapes.json> [isolation] [txns]}"
iso="${2:-read-committed}"
txns="${3:-300}"

command -v elle-go >/dev/null || { echo "elle-go not found on PATH; run 'go install .' in the elle-go repo" >&2; exit 2; }
docker info >/dev/null 2>&1 || { echo "Docker is not running" >&2; exit 2; }

name="elle-check-$$"
docker run -d --rm --name "$name" -e POSTGRES_PASSWORD=elle -p 127.0.0.1::5432 postgres:16 >/dev/null
trap 'docker stop "$name" >/dev/null 2>&1 || true' EXIT

for _ in $(seq 1 120); do
  docker exec "$name" pg_isready -h 127.0.0.1 -U postgres >/dev/null 2>&1 && break
  sleep 0.5
done
port="$(docker port "$name" 5432/tcp | head -1 | sed 's/.*://')"

set +e
elle-go pg -dsn "postgres://postgres:elle@127.0.0.1:${port}/postgres?sslmode=disable" \
  -isolation "$iso" -shapes "$shapes" -txns "$txns"
code=$?
set -e
exit "$code"
