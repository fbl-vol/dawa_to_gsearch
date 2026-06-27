#!/usr/bin/env sh
set -eu

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for this example" >&2
  exit 1
fi

API_KEY="${DATAFORDELEREN_API_KEY:-${DATAFORDELER_GRAPHQL_API_KEY:-}}"

if [ -z "$API_KEY" ]; then
  echo "Set DATAFORDELEREN_API_KEY or DATAFORDELER_GRAPHQL_API_KEY" >&2
  exit 1
fi

EASTING="${EASTING:-689255}"
NORTHING="${NORTHING:-6051787}"
BBOX_SIZE_METERS="${BBOX_SIZE_METERS:-20}"
LIMIT="${LIMIT:-10}"
NOW="${NOW:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

HALF=$((BBOX_SIZE_METERS / 2))
MIN_E=$((EASTING - HALF))
MIN_N=$((NORTHING - HALF))
MAX_E=$((EASTING + HALF))
MAX_N=$((NORTHING + HALF))

WKT="POLYGON(($MIN_E $MIN_N, $MAX_E $MIN_N, $MAX_E $MAX_N, $MIN_E $MAX_N, $MIN_E $MIN_N))"

QUERY='query GetBygninger($first: Int, $where: GEODKV_BygningFilterInput, $registreringstid: DafDateTime, $virkningstid: DafDateTime) {
  GEODKV_Bygning(first: $first, where: $where, registreringstid: $registreringstid, virkningstid: $virkningstid) {
    nodes {
      geometri { crs wkt }
      bygningstype
      status
      id_lokalId
      BBRUUID
    }
  }
}'

VARIABLES=$(jq -n \
  --argjson first "$LIMIT" \
  --arg now "$NOW" \
  --arg wkt "$WKT" \
  '{
    first: $first,
    registreringstid: $now,
    virkningstid: $now,
    where: {
      geometri: {
        intersects: {
          crs: 25832,
          wkt: $wkt
        }
      }
    }
  }')

BODY=$(jq -n \
  --arg query "$QUERY" \
  --argjson variables "$VARIABLES" \
  '{query: $query, variables: $variables}')

ENCODED_KEY=$(printf '%s' "$API_KEY" | jq -sRr @uri)

curl -sS -X POST \
  -H 'Content-Type: application/json' \
  "https://graphql.datafordeler.dk/GEODKV/v2?apiKey=$ENCODED_KEY" \
  --data "$BODY"
