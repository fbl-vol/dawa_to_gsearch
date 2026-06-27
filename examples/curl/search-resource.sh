#!/usr/bin/env sh
set -eu

: "${GSEARCH_TOKEN:?Set GSEARCH_TOKEN to a Dataforsyningen token}"

BASE_URL="${GSEARCH_BASE_URL:-https://api.dataforsyningen.dk/rest/gsearch/v2.0}"
RESOURCE="${1:-husnummer}"
QUERY="${2:-}"
LIMIT="${LIMIT:-5}"

case "$RESOURCE" in
  adresse) default_query="flens" ;;
  husnummer) default_query="genvej" ;;
  navngivenvej) default_query="krin" ;;
  postnummer) default_query="mari" ;;
  kommune) default_query="aalborg" ;;
  region) default_query="midt" ;;
  sogn) default_query="budolfi" ;;
  politikreds) default_query="vest" ;;
  retskreds) default_query="københavn" ;;
  opstillingskreds) default_query="vest" ;;
  matrikel) default_query="123ab" ;;
  matrikel_udgaaet) default_query="11a" ;;
  stednavn) default_query="Benedikte" ;;
  *)
    echo "Unsupported resource: $RESOURCE" >&2
    echo "Use one of: adresse husnummer navngivenvej postnummer kommune region sogn politikreds retskreds opstillingskreds matrikel matrikel_udgaaet stednavn" >&2
    exit 64
    ;;
esac

if [ -z "$QUERY" ]; then
  QUERY="$default_query"
fi

curl -sS --get "$BASE_URL/$RESOURCE" \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode "q=$QUERY" \
  --data-urlencode "limit=$LIMIT"
