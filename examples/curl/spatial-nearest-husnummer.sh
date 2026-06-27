#!/usr/bin/env sh
set -eu

: "${GSEARCH_TOKEN:?Set GSEARCH_TOKEN to a Dataforsyningen token}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for this example" >&2
  exit 69
fi

# EPSG:25832 point near Genvej 1, 4874 Gedser. Convert WGS84 to EPSG:25832 in
# your application before calling GSearch spatial filters.
EASTING="${1:-689255}"
NORTHING="${2:-6051787}"
POSTAL_PREFIX="${3:-4}"

point="POINT($EASTING $NORTHING)"

postnummer_json=$(
  curl -sS --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/postnummer' \
    --data-urlencode "token=$GSEARCH_TOKEN" \
    --data-urlencode "q=$POSTAL_PREFIX" \
    --data-urlencode 'limit=5' \
    --data-urlencode "filter=INTERSECTS(geometri,$point)"
)

postnummernavn=$(printf '%s' "$postnummer_json" | jq -r '.[0].postnummernavn // empty')
if [ -z "$postnummernavn" ]; then
  echo "No containing postnummer found for $point" >&2
  exit 1
fi

for radius in 100 300 1000; do
  husnummer_json=$(
    curl -sS --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer' \
      --data-urlencode "token=$GSEARCH_TOKEN" \
      --data-urlencode "q=$postnummernavn" \
      --data-urlencode 'limit=5' \
      --data-urlencode "filter=DWITHIN(geometri,$point,$radius,meters)"
  )

  if [ "$(printf '%s' "$husnummer_json" | jq 'length')" -gt 0 ]; then
    printf '%s' "$husnummer_json" | jq --arg radius "$radius" --arg postnummernavn "$postnummernavn" '{
      postnummernavn: $postnummernavn,
      radius_meters: ($radius | tonumber),
      candidates: [.[] | {
        id,
        visningstekst,
        kommunekode,
        vejkode,
        postnummer,
        geometri
      }]
    }'
    exit 0
  fi
done

echo "No nearby husnummer found within 1000 meters of $point" >&2
exit 1
