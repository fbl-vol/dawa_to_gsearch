#!/usr/bin/env sh
set -eu

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for this example" >&2
  exit 69
fi

API_KEY="${DATAFORDELER_GRAPHQL_API_KEY:-${DATAFORDELEREN_API_KEY:-}}"

if [ -z "$API_KEY" ]; then
  echo "Set DATAFORDELER_GRAPHQL_API_KEY or DATAFORDELEREN_API_KEY" >&2
  exit 1
fi

REGISTER="${1:-DAR}"
QUERY="${2:-query { __typename }}"

case "$REGISTER" in
  DAR|BBR|MAT) ;;
  *)
    echo "Unsupported register: $REGISTER" >&2
    echo "Use one of: DAR BBR MAT" >&2
    exit 64
    ;;
esac

key_param="apiKey"
encoded_key=$(printf '%s' "$API_KEY" | jq -sRr @uri)
url="https://graphql.datafordeler.dk/$REGISTER/v2?$key_param=$encoded_key"
payload=$(jq -n --arg query "$QUERY" '{query: $query}')

curl -sS "$url" \
  -H 'accept: application/json' \
  -H 'content-type: application/json' \
  --data "$payload"
