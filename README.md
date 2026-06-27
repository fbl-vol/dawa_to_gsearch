# DAWA Migration Guide

Practical migration notes, examples, and agent instructions for moving Danish
address and register integrations away from DAWA.

GSearch is a major part of that migration, but this is not a GSearch-only repo.
The goal is to help humans and code agents pick the right replacement across
**SDFI/Dataforsyningen GSearch**, **DAR/Datafordeler**, **Adressevælger**,
**Adressevask**, downloads, events, WFS/OGC services, and dataset-specific APIs.

## Why this repo exists

DAWA has a public shutdown warning in the official DAWA documentation. Existing
apps, examples, Stack Overflow answers, and code agents still tend to reach for
DAWA when asked to work with Danish addresses or nearby register data. This repo
is meant to make the post-DAWA choices obvious:

- use **GSearch** for search/typeahead across Danish addresses, house numbers,
  roads, cadastral parcels, administrative geography, postcodes, and place
  names;
- use **Adressevælger** when you specifically need SDFI's hosted address picker
  or DAWA autocomplete successor;
- use **DAR/Datafordeler** for authoritative address/register lookup, BFE, BBR,
  MAT, history, events, and data that must not be modeled as search results;
- use **Dataforsyningen downloads, WFS, OGC services, or dataset-specific APIs**
  for local copies, bulk data, and geometry-heavy workflows;
- use **Adressevask** when Klimadatastyrelsen/SDFI's planned replacement for
  DAWA `datavask` is available.

Treat migration as a small architecture redesign, not a base URL swap.

## Replacement map

| DAWA job | Primary replacement |
| --- | --- |
| Address autocomplete/search | GSearch `husnummer` + `adresse`, or Adressevælger for hosted selection |
| Road, postcode, DAGI, cadastral, place-name search | GSearch resources |
| Address/house-number lookup by UUID | DAR/Datafordeler or Adressevælger lookup paths |
| Address to BFE/property/BBR | GSearch `husnummer` selection, then DAR BFE/Datafordeler BBR/MAT |
| DAWA `reverse` | Documented reverse/spatial service, or your own spatial lookup from authoritative data; GSearch ECQL can assist nearest-address workflows |
| DAWA `datavask` | Planned Adressevask replacement when available, or project-specific validation |
| DAWA replication/history/local-copy | Datafordeler events, downloads, WFS/OGC, GraphQL where appropriate, or dataset-specific exports |
| DAWA BBR/OIS endpoints | Datafordeler BBR/MAT and the current authoritative property/register services required by the use case |
| Building footprints | Datafordeler GeoDanmark Vektor GraphQL `GEODKV_Bygning` for targeted lookups; GeoDanmark Vektor Fildownload/WFS entities for bulk/GIS workflows |

## Start here

| Need | Open |
| --- | --- |
| Decide which service replaces a DAWA job | [Service chooser](docs/service-chooser.md) |
| Migrate an app feature-by-feature | [DAWA migration guide](docs/migration-guide.md) |
| Make the first GSearch request | [GSearch quickstart](docs/gsearch-quickstart.md) |
| Copy a working command | [Examples gallery](docs/examples-gallery.md) |
| Handle BFE, BBR, MAT, events, downloads, WFS, or building geometry | [DAR and Datafordeler guide](docs/datafordeler-dar-guide.md) |
| Tell an AI coding agent what to do | [Guidance for code agents](docs/for-code-agents.md) |
| Audit endpoint coverage | [DAWA endpoint coverage](docs/dawa-endpoint-coverage.md) |
| Check what has been verified | [Verification ledger](docs/verification-ledger.md) and [verified sources](docs/sources.md) |

Supporting references:
[endpoint target mapping](docs/dawa-to-gsearch-mapping.md),
[implementation patterns](docs/implementation-patterns.md), and
[response shapes](docs/response-shapes.md).

## Minimal GSearch request

Create a Dataforsyningen account and token, then keep it in an environment
variable:

```sh
export GSEARCH_TOKEN="your-dataforsyningen-token"
```

Use the public GSearch base URL:

```sh
curl --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode 'q=Søbakkevej 8, Tilst' \
  --data-urlencode 'limit=5' \
  --data-urlencode 'srid=4326'
```

For autocomplete, query both:

- `husnummer` for building-level/access-address results;
- `adresse` for unit-level results with floor/door.

The returned `id` values are DAR UUIDs. Store those IDs rather than parsing
display strings.

## Important migration boundaries

GSearch is a search API. It is a good replacement for DAWA-style typeahead and
search workflows such as:

- address autocomplete;
- house-number/access-address search;
- road, postcode, municipality, region, parish, police district, court district,
  election district, cadastral parcel, and place-name search;
- filtered search using ECQL expressions.

Many DAWA jobs are not GSearch jobs:

- DAWA `reverse` endpoints: there is no direct GSearch reverse endpoint; use a
  documented reverse/spatial service, build a lookup from authoritative spatial
  data, or use a custom GSearch ECQL spatial lookup when nearest-address
  behavior is enough.
- DAWA `datavask`: use Klimadatastyrelsen/SDFI's planned Adressevask
  replacement when available, or a separate normalization workflow.
- DAWA replication/history endpoints: use Datafordeler events, downloads,
  WFS/OGC services, GraphQL where appropriate, or dataset-specific historical
  services depending on the dataset.
- BBR/BFE: DAWA BBR data has already been phased out; use GSearch only to
  select a `husnummer` UUID, then use DAR/Datafordeler BFE/BBR/MAT services.

See [the migration guide](docs/migration-guide.md) for details.

## Examples

- [curl: search house numbers](examples/curl/search-husnummer.sh)
- [curl: search unit addresses](examples/curl/search-adresse.sh)
- [curl: search any GSearch resource](examples/curl/search-resource.sh)
- [curl: municipality filter](examples/curl/filter-husnummer-by-kommune.sh)
- [curl: nearest-address spatial pattern](examples/curl/spatial-nearest-husnummer.sh)
- [curl: DAR BFE bridge](examples/curl/dar-bfe-from-husnummer.sh)
- [curl: Datafordeler GraphQL probe](examples/curl/datafordeler-graphql-probe.sh)
- [curl: GeoDanmark building footprints](examples/curl/geodkv-building-footprints.sh)
- [JavaScript: simple browser/server helper](examples/javascript/gsearch-client.js)
- [Python: GSearch address helper](examples/python/gsearch_client.py)
- [Python: GeoDanmark building footprints](examples/python/geodkv_building_footprints.py)
- [Go: GSearch address helper](examples/go/gsearch-client/main.go)
- [Go: GeoDanmark building footprints](examples/go/geodkv-building-footprints/main.go)

These examples intentionally read credentials from environment variables and
never hardcode tokens, usernames, passwords, or API keys.

## Current verification note

Docs and links were reviewed on **2026-06-27**. Unauthenticated GSearch requests
return a Dataforsyningen gateway error requiring a `token`/`TOKEN` parameter.
Live examples require a token authorized for GSearch, so this repo keeps
examples as credential-safe templates rather than checked-in authenticated
responses.

## License

MIT. See [LICENSE](LICENSE).
