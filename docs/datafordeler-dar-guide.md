# DAR and Datafordeler guide

Use this page when the old DAWA call was doing more than user-facing search.
The post-DAWA pattern is usually:

1. Let the user select an address or object with GSearch or Adressevælger.
2. Persist the returned UUIDs and structured fields.
3. Use DAR/Datafordeler, BBR, MAT, events, downloads, WFS/OGC, or another
   authoritative service for the data workflow.

## Authentication shapes

Different replacement services use different authentication paths.

| Service family | Typical credential shape | Example in this repo |
| --- | --- | --- |
| GSearch | Dataforsyningen `token` request parameter | `GSEARCH_TOKEN` |
| Datafordeler legacy REST | `username` + `password` parameters | `DATAFORDELER_USERNAME`, `DATAFORDELER_PASSWORD` |
| Datafordeler GraphQL | API key or OAuth through Datafordeler Administration | `DATAFORDELER_GRAPHQL_API_KEY`; `DATAFORDELEREN_API_KEY` is also accepted by GraphQL examples as a common alias |
| Adressevælger/Adressevask | Token parameter according to the official FAQ; status can change with brugerstyring | Follow Klimadatastyrelsen docs |

Do not collapse these into one generic `API_KEY`. Keep service-specific
configuration names so a future maintainer can see which platform a credential
belongs to.

## Address to BFE

When an old DAWA flow used an address to reach property or building data, the
important bridge is usually the `husnummer` UUID.

Flow:

1. Search `husnummer` and `adresse` with GSearch, or select an address with
   Adressevælger.
2. Keep the parent `husnummer` UUID.
3. Resolve that UUID through DAR BFE to get `jordstykkeLokalId` and
   `samletFastEjendom` (BFE number).
4. Use those values with BBR, MAT, or property-specific Datafordeler workflows.

Example:

```sh
examples/curl/dar-bfe-from-husnummer.sh "<husnummer-uuid>"
```

The DAR BFE REST page is marked for phase-out by Datafordeler, so use this as a
bridge for migration and verify the current GraphQL/API-key/OAuth path before
building new long-lived integrations.

## DAR GraphQL

Use DAR GraphQL when the application needs authoritative address entities rather
than search results. Official Datafordeler docs describe entities such as
Adresse, Adressepunkt, Husnummer, NavngivenVej, NavngivenVejKommunedel,
NavngivenVejPostnummer, NavngivenVejSupplerendeBynavn, Postnummer, and
SupplerendeBynavn.

Probe access:

```sh
examples/curl/datafordeler-graphql-probe.sh DAR
```

The default probe query is deliberately tiny:

```graphql
query { __typename }
```

Real queries should request the entity fields your application actually needs.
Do not copy DAWA response shapes blindly into GraphQL.

## BBR and MAT

Use BBR GraphQL, BBR REST during transition, MAT GraphQL, or relevant
Datafordeler services when the old DAWA flow touched:

- buildings;
- units;
- floors;
- entrances;
- ground/property data;
- technical installations;
- BFE/property/cadastral relations.

Probe BBR or MAT GraphQL access:

```sh
examples/curl/datafordeler-graphql-probe.sh BBR
examples/curl/datafordeler-graphql-probe.sh MAT
```

BBR REST is also marked for phase-out by Datafordeler. Prefer current
Datafordeler GraphQL guidance for new integrations.

## GeoDanmark building geometry

DAWA `bygninger` exposed building polygons from GeoDanmark data. Do not replace
that with GSearch. For targeted building-footprint lookup, use Datafordeler
GeoDanmark Vektor GraphQL:

- register: `GEODKV`;
- version: `v2`;
- endpoint: `https://graphql.datafordeler.dk/GEODKV/v2`;
- entity: `GEODKV_Bygning`;
- geometry field: `geometri`;
- useful fields: `geometri { crs wkt }`, `bygningstype`, `status`,
  `id_lokalId`, and `BBRUUID`;
- authentication: API key or OAuth through Datafordeler Administration.

Use `GEODKV_BygningFilterInput.geometri.intersects` with EPSG:25832 WKT when
you need buildings near an address/access point. Use GeoDanmark Vektor
Fildownload or WFS entities when the workflow needs bulk data, a local copy, or
GIS tooling.

Example:

```sh
examples/curl/geodkv-building-footprints.sh
```

Legacy GeoDanmark WFS and FTP-style `filudtræk` pages have phase-out notices.
Check the current Datafordeler data overview before starting new long-lived
building-geometry integrations.

## Events and local copies

DAWA replication endpoints were local-copy/change-feed workflows. Replace them
with Datafordeler events, file downloads, WFS entity services, OGC services, or
dataset-specific exports.

The Datafordeler DAR and BBR event pages describe register-specific event
entities such as `DAR_Events` and `BBR_Events`. Use those when your system needs
to keep its own copy synchronized.

For full refreshes or local analytical stores, start with Datafordeler
fildownload/WFS or Dataforsyningen downloads rather than repeated GSearch
queries.

## Reverse geocoding

For authoritative reverse-geocoding behavior, use a documented service for that
specific job, or build your own spatial index from DAR/Datafordeler, WFS,
downloads, or another authoritative dataset. GSearch can assist nearest-address
product workflows with ECQL spatial filters, but it is not a replacement for
DAWA reverse endpoint semantics.

See [implementation patterns](implementation-patterns.md) for the
`INTERSECTS`/`DWITHIN` pattern.

## What to keep from GSearch

Even in Datafordeler-heavy systems, GSearch remains useful at the edge:

- human search text;
- human-entered text search and typeahead matching;
- `visningstekst` for suggestion display;
- structured fields such as `kommunekode`, `vejkode`, `postnummer`;
- the selected `adresse` and `husnummer` UUIDs.

After selection, switch to the authoritative service for the actual data job.
