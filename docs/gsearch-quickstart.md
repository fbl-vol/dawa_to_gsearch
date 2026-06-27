# GSearch quickstart

GSearch is SDFI/Dataforsyningen's REST search API for Danish geodata. It is
well-suited to autocomplete/typeahead and "find the best matching official
object" workflows.

It is not Google Search, and it is not a full clone of DAWA.

## 1. Get a token

Create a user at <https://dataforsyningen.dk/> and create a token.

Use an environment variable in local development:

```sh
export GSEARCH_TOKEN="your-dataforsyningen-token"
```

The general Dataforsyningen webservice documentation shows token usage as a
request parameter named `token`:

```txt
https://api.dataforsyningen.dk/<service>?<parameters>&token=<token>
```

Some existing internal clients pass a lowercase `token` header. Verify that
works in your environment before relying on it. For public examples, prefer the
query parameter because it matches the platform documentation and the current
gateway error messages.

## 2. Pick the right resource

Base URL:

```txt
https://api.dataforsyningen.dk/rest/gsearch/v2.0/{resource}
```

Common resources:

| Resource | Use it for |
| --- | --- |
| `adresse` | Specific addresses/units, including floor and door where present |
| `husnummer` | Building-level access addresses/house numbers |
| `navngivenvej` | Named roads |
| `postnummer` | Postcodes/postal districts |
| `kommune` | Municipalities |
| `region` | Regions |
| `sogn` | Parishes |
| `politikreds` | Police districts |
| `retskreds` | Court districts |
| `opstillingskreds` | Election nomination districts |
| `matrikel` | Current cadastral parcels |
| `matrikel_udgaaet` | Retired cadastral parcels |
| `stednavn` | Danish place names |

For address autocomplete, query both `husnummer` and `adresse`. Use
`husnummer` when the user is selecting a building/access address, and `adresse`
when the user may select a unit with floor/door.

To try all resources with verified sample queries, see
[examples gallery](examples-gallery.md).

## 3. Make a request

House-number search:

```sh
curl --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode 'q=Søbakkevej 8, Tilst' \
  --data-urlencode 'limit=5' \
  --data-urlencode 'srid=4326'
```

Unit/address search:

```sh
curl --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/adresse' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode 'q=Søbakkevej 8, Tilst' \
  --data-urlencode 'limit=5' \
  --data-urlencode 'srid=4326'
```

Road search:

```sh
curl --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/navngivenvej' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode 'q=Lærke' \
  --data-urlencode 'limit=10'
```

Generic resource search:

```sh
examples/curl/search-resource.sh kommune aalborg
examples/curl/search-resource.sh matrikel 123ab
examples/curl/search-resource.sh stednavn Benedikte
```

## 4. Understand the important parameters

| Parameter | Required | Notes |
| --- | --- | --- |
| `token` | Yes | Dataforsyningen token. The gateway requires the parameter name to be `token` or `TOKEN`. |
| `q` | Yes | Search text. GSearch handles common spelling, writing, and phonetic variants. |
| `limit` | No | Maximum response count. Default is 10; maximum is 100. |
| `filter` | No | ECQL filter expression. Fully URL-encode it. |
| `srid` | No | Returned geometry coordinate system. Use `4326` for longitude/latitude when supported by the endpoint. |

GSearch treats `-`, `(`, `)`, and `!` in `q` as spaces.

## 5. Filters

Use `filter` to constrain a search by attributes or geometry. GSearch filters
use ECQL, the GeoServer extension of OGC CQL.

Municipality filter example:

```sh
curl --get 'https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer' \
  --data-urlencode "token=$GSEARCH_TOKEN" \
  --data-urlencode 'q=lærke' \
  --data-urlencode "filter=kommunekode like '%0461%'" \
  --data-urlencode 'limit=20'
```

Spatial filters must use EPSG:25832 coordinates, even if you request returned
geometry as `srid=4326`.

```txt
INTERSECTS(geometri,POLYGON((515000 6074200,515000 6104200,555000 6104200,555000 6074200,515000 6074200)))
```

For `adresse` and `husnummer`, use `geometri` or `vejpunkt_geometri`. For
`matrikel` and `matrikel_udgaaet`, `geometri`, `centroid_x`, and `centroid_y`
are relevant.

Spatial filters can support custom map workflows, including nearest-address
lookup with `DWITHIN`, but they are not DAWA `/reverse` endpoint parity. See
[implementation patterns](implementation-patterns.md) before migrating reverse
geocoding.

## 6. Response basics

Every result includes:

- `id`: authoritative UUID for the matched object;
- `visningstekst`: display text for suggestion lists;
- geometry fields, usually GeoJSON;
- resource-specific fields such as `kommunekode`, `vejkode`, `postnummer`, and
  road/address components.

For migration work, persist the UUID and structured fields. Do not parse
`visningstekst` as your source of truth.

Do not treat GSearch path lookups such as `/adresse/{id}` or
`/husnummer/{id}` as DAWA parity. Live verification on 2026-06-27 returned 404
for those path forms. Use DAR/Datafordeler or Adressevælger lookup paths when
your application needs durable lookup by UUID.

## Troubleshooting

| Symptom | Likely cause | Fix |
| --- | --- | --- |
| `400 Authentication token must only be lowercase (token) or uppercase (TOKEN)!` | Missing token parameter, wrong parameter name, or a rejected token header path | Send `token=<token>` or `TOKEN=<token>` as a request parameter. |
| `401 User not authorized` | Token exists but is not authorized for this service | Create/check the token in Dataforsyningen and verify service access. |
| Empty results for valid-looking text | Resource mismatch or overly narrow filter | Query both `husnummer` and `adresse`, remove filters, then add filters back. |
| Wrong coordinates | SRID mismatch | Use `srid=4326` for WGS84 output, and EPSG:25832 for spatial filters. |
| Slow or timed-out request | Broad search or expensive geometry filter | Add a tighter `q`, reduce `limit`, or simplify `filter`. GSearch docs note that requests over 10 seconds can return 504. |

## Minimal autocomplete algorithm

1. Debounce input for 150-300 ms.
2. Skip requests until the user has typed at least 2-3 meaningful characters.
3. Query `husnummer` and `adresse` in parallel.
4. Merge results by structured fields, not only display text.
5. Prefer `adresse` for unit-level display; retain the parent `husnummer` ID
   when you need building-level lookups later.
6. Store the selected result's `id`, type, display text, and structured fields.
