#!/usr/bin/env sh
set -eu

: "${GSEARCH_TOKEN:?Set GSEARCH_TOKEN to a Dataforsyningen token}"

QUERY="${1:-lærke}"
KOMMUNEKODE="${2:-0461}"
LIMIT="${LIMIT:-20}"

curl -sS --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode "q=$QUERY" \
  --data-urlencode "filter=kommunekode like '%$KOMMUNEKODE%'" \
  --data-urlencode "limit=$LIMIT"
