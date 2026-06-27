#!/usr/bin/env sh
set -eu

: "${GSEARCH_TOKEN:?Set GSEARCH_TOKEN to a Dataforsyningen token}"

QUERY="${1:-Søbakkevej 8, Tilst}"
LIMIT="${LIMIT:-5}"

curl -sS --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode "q=$QUERY" \
  --data-urlencode "limit=$LIMIT" \
  --data-urlencode 'srid=4326'
