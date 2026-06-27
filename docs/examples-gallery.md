# Examples gallery

This page gives one working shape for each major migration job: GSearch search,
filters, custom spatial lookup, DAR BFE bridging, GraphQL probing, building
geometry, and the non-GSearch DAWA features.

All examples expect secrets in environment variables. Do not paste real tokens
into source code.

## Setup

Set these environment variables in your shell or secret manager:

```txt
GSEARCH_TOKEN
DATAFORDELER_USERNAME
DATAFORDELER_PASSWORD
DATAFORDELER_GRAPHQL_API_KEY
DATAFORDELEREN_API_KEY
```

`DATAFORDELER_GRAPHQL_API_KEY` and `DATAFORDELEREN_API_KEY` are both used in
real codebases for Datafordeler GraphQL keys. The GraphQL examples accept either
name.

The Python examples use only the standard library and target Python 3.9+. The
Go examples use only the standard library and target Go 1.18+.

## GSearch resources

Use the generic resource example to try every first-class GSearch resource:

```sh
examples/curl/search-resource.sh husnummer
examples/curl/search-resource.sh adresse
examples/curl/search-resource.sh navngivenvej
examples/curl/search-resource.sh postnummer
examples/curl/search-resource.sh kommune
examples/curl/search-resource.sh region
examples/curl/search-resource.sh sogn
examples/curl/search-resource.sh politikreds
examples/curl/search-resource.sh retskreds
examples/curl/search-resource.sh opstillingskreds
examples/curl/search-resource.sh matrikel
examples/curl/search-resource.sh matrikel_udgaaet
examples/curl/search-resource.sh stednavn
```

The defaults above were live-probed on 2026-06-27. Each returned HTTP 200 with
an authorized GSearch token.

| Resource | Default query | First verified result |
| --- | --- | --- |
| `husnummer` | `genvej` | `Genvej 1, 4874 Gedser` |
| `adresse` | `flens` | `Flensbjerg 1A, 4960 Holeby` |
| `navngivenvej` | `krin` | `Kringelborg Alle (4800 Nykøbing F)` |
| `postnummer` | `mari` | `9550 Mariager` |
| `kommune` | `aalborg` | `Aalborg Kommune` |
| `region` | `midt` | `Region Midtjylland` |
| `sogn` | `budolfi` | `Budolfi sogn` |
| `politikreds` | `vest` | `Københavns Vestegns Politikreds` |
| `retskreds` | `københavn` | `Københavns Byret` |
| `opstillingskreds` | `vest` | `Vesterbrokredsen` |
| `matrikel` | `123ab` | `123ab, Povlsker` |
| `matrikel_udgaaet` | `11a` | `11ae, Em By, Em` |
| `stednavn` | `Benedikte` | `Benedikte Sø (Sø i Gråsten)` |

Use explicit queries when you want a different example:

```sh
examples/curl/search-resource.sh husnummer "Søbakkevej 8, Tilst"
examples/curl/search-resource.sh adresse "Søbakkevej 8, Tilst"
LIMIT=20 examples/curl/search-resource.sh postnummer aar
```

## Address autocomplete

For a custom address picker, query both address resources:

```sh
examples/curl/search-husnummer.sh "Søbakkevej 8, Tilst"
examples/curl/search-adresse.sh "Søbakkevej 8, Tilst"
```

In application code, merge results by structured address fields. Prefer
`adresse` for unit-level display, but keep the parent `husnummer` UUID for
BBR/BFE/property workflows.

The JavaScript helper shows that merge:

```js
import { searchAddressSuggestions } from "./examples/javascript/gsearch-client.js";

const suggestions = await searchAddressSuggestions("Søbakkevej 8, Tilst");
console.log(suggestions[0]);
```

Python, using only the standard library:

```sh
python3 examples/python/gsearch_client.py "Søbakkevej 8, Tilst"
```

Go, using only the standard library:

```sh
go run examples/go/gsearch-client/main.go "Søbakkevej 8, Tilst"
```

Both examples query `husnummer` and `adresse` in parallel, merge by structured
address fields, preserve `husnummerId`, and handle GeoJSON `Point`, `MultiPoint`,
and nested geometry coordinate shapes.

## Attribute filters

GSearch filters use ECQL. This example searches house numbers and limits results
to one municipality code:

```sh
examples/curl/filter-husnummer-by-kommune.sh lærke 0461
```

The script sends:

```txt
filter=kommunekode like '%0461%'
```

## Spatial lookup

GSearch has no DAWA `/reverse` endpoint, but ECQL spatial filters can support a
nearest-address workflow when endpoint parity is not required.

This example uses an EPSG:25832 point near `Genvej 1, 4874 Gedser`:

```sh
examples/curl/spatial-nearest-husnummer.sh
```

It first finds the containing `postnummer`:

```txt
INTERSECTS(geometri,POINT(689255 6051787))
```

Then it searches nearby `husnummer` candidates:

```txt
DWITHIN(geometri,POINT(689255 6051787),100,meters)
```

Use a real coordinate transform in production. Convert WGS84 latitude/longitude
to EPSG:25832 before building the filter.

## Address to BFE/BBR

Use GSearch to select the `husnummer` UUID. Then resolve that UUID through DAR
BFE:

```sh
examples/curl/dar-bfe-from-husnummer.sh "0a3f5087-05a4-32b8-e044-0003ba298018"
```

The response shape contains `jordstykkeList`; when a match exists, use
`jordstykkeLokalId` and `samletFastEjendom` as bridge values into BBR/MAT or
other property workflows.

DAR BFE REST is marked for phase-out by Datafordeler, so use this as a migration
bridge and check the current GraphQL/API-key/OAuth path before building new
long-lived integrations.

## Datafordeler GraphQL

Probe Datafordeler GraphQL access without exposing the API key:

```sh
examples/curl/datafordeler-graphql-probe.sh DAR
examples/curl/datafordeler-graphql-probe.sh BBR
examples/curl/datafordeler-graphql-probe.sh MAT
```

The default query is:

```graphql
query { __typename }
```

For real integrations, write entity-specific queries from the Datafordeler DAR,
BBR, MAT, and events documentation. The exact subfields should be driven by the
data your app needs, not by old DAWA response shapes.

See [DAR and Datafordeler guide](datafordeler-dar-guide.md) for the broader
authoritative-data migration pattern.

## Adressevælger

Use Adressevælger when you want Klimadatastyrelsen/SDFI's official hosted
address selector or API path for the DAWA autocomplete value proposition.

The official FAQ says it returns DAR UUIDs, is not a 1:1 response/interface
replacement for DAWA autocomplete, and does not provide DAWA-style reverse
geocoding. Treat an Adressevælger migration as a UI/API adaptation, not a URL
swap.

## Adressevask

Use the planned Adressevask replacement for DAWA `datavask` only when it is
available. As of the official page reviewed on 2026-06-27, Adressevask is not
available yet and is expected at the end of August 2026.

Until then, do not add new DAWA `datavask` dependencies. If an existing
environment still has a temporary DAWA `datavask` dependency, treat it as
migration debt, or build a project-specific normalization/validation flow around
the official data sources you need.

## Replication, history, and local copies

Do not try to make GSearch a replication feed. Use Datafordeler events,
Datafordeler file downloads, WFS, OGC services, or dataset-specific exports.

For event-driven replacement work, start with the register's event service, for
example DAR events or BBR events. For complete local copies, start with
Datafordeler file downloads or WFS entity services.

## Building footprints

Do not use GSearch as a building-footprint source. Use GSearch only to resolve
the user's address selection to coordinates and IDs, then switch to Datafordeler
GeoDanmark Vektor.

For targeted building geometry near a selected address, use GeoDanmark Vektor
GraphQL `GEODKV_Bygning`:

```sh
examples/curl/geodkv-building-footprints.sh
```

Python:

```sh
python3 examples/python/geodkv_building_footprints.py
```

Go:

```sh
go run examples/go/geodkv-building-footprints/main.go
```

The example queries `https://graphql.datafordeler.dk/GEODKV/v2` with a spatial
`intersects` filter and asks for `geometri { crs wkt }`, `bygningstype`,
`status`, `id_lokalId`, and `BBRUUID`.

GraphQL examples expect either:

```txt
DATAFORDELEREN_API_KEY
DATAFORDELER_GRAPHQL_API_KEY
```

For bulk or local-copy workflows, use GeoDanmark Vektor Fildownload or WFS
entities. Avoid new dependencies on legacy GeoDanmark WFS/FTP filudtræk paths
without checking the current Datafordeler data overview and phase-out notices.
