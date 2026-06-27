#!/usr/bin/env sh
set -eu

: "${DATAFORDELER_USERNAME:?Set DATAFORDELER_USERNAME}"
: "${DATAFORDELER_PASSWORD:?Set DATAFORDELER_PASSWORD}"

HUSNUMMER_ID="${1:?Usage: DATAFORDELER_USERNAME=... DATAFORDELER_PASSWORD=... $0 <husnummer-uuid>}"

curl -sS --get 'https://services.datafordeler.dk/DAR/DAR_BFE_Public/1/REST/husnummerTilBygningBfe' \
  --data-urlencode "HusnummerId=$HUSNUMMER_ID" \
  --data-urlencode 'format=json' \
  --data-urlencode "username=$DATAFORDELER_USERNAME" \
  --data-urlencode "password=${DATAFORDELER_PASSWORD}"
