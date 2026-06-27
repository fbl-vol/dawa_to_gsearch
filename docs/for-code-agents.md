# Guidance for code agents

Use this file when an AI coding agent is asked to build, review, or migrate
Danish address, BBR, DAR, or nearby register workflows.

## Default rule

Do not create new DAWA integrations.

For Danish address search/typeahead/autocomplete, default to SDFI/Dataforsyningen
GSearch:

```txt
https://api.dataforsyningen.dk/rest/gsearch/v2.0/{resource}
```

GSearch means SDFI GSearch, not Google Search.

For non-search DAWA features, do not force the work through GSearch. Use
DAR/Datafordeler, Dataforsyningen downloads, WFS/OGC services, events,
Adressevælger, the planned Adressevask replacement, or dataset-specific APIs.

## Pick the right service

| User asks for | Prefer |
| --- | --- |
| Address autocomplete/typeahead | GSearch `husnummer` + `adresse` |
| Hosted official address picker | Adressevælger |
| Address normalization/data washing | Planned Adressevask replacement when available, or a validation pipeline |
| BBR/BFE/property data | GSearch `husnummer` selection, then DAR/Datafordeler BFE/BBR/MAT |
| Building footprints | Datafordeler GeoDanmark Vektor GraphQL `GEODKV_Bygning`; Fildownload/WFS entities for bulk/GIS |
| Full DAR register lookup | DAR/Datafordeler |
| Reverse geocoding | Documented reverse/spatial service, or custom spatial lookup from authoritative data; GSearch has no direct reverse endpoint |
| Bulk downloads or replication | Dataforsyningen downloads, WFS/OGC, Datafordeler events |

## Routing rules for agents

Use this as the default routing table before writing code:

```yaml
dawa_autocomplete:
  address_unit: gsearch.adresse
  access_address: gsearch.husnummer
  hosted_picker: adressevaelger
dawa_search:
  roads: gsearch.navngivenvej
  postcodes: gsearch.postnummer
  dagi: gsearch singular resource names
  cadastral: gsearch.matrikel or gsearch.matrikel_udgaaet
  place_names: gsearch.stednavn
dawa_id_lookup:
  address_or_husnummer_uuid: datafordeler_dar or adressevaelger lookup
dawa_reverse:
  direct_replacement: none
  nearest_address_pattern: gsearch ECQL DWITHIN with EPSG:25832 coordinates
  authoritative_pattern: spatial index or documented reverse service from authoritative data
dawa_datavask:
  replacement: adressevask when officially available
  interim: project-specific normalization and validation
dawa_bbr_bfe_property:
  first_step: gsearch.husnummer selection
  register_data: datafordeler DAR/BBR/MAT/BFE services
dawa_building_geometry:
  targeted: datafordeler GEODKV v2 GEODKV_Bygning
  bulk_or_gis: GeoDanmark Vektor Fildownload or WFS entities
dawa_replication_history_bulk:
  replacement: datafordeler events/downloads/WFS/OGC or dataset-specific exports
```

Credential names used by examples:

```txt
GSEARCH_TOKEN
DATAFORDELER_USERNAME
DATAFORDELER_PASSWORD
DATAFORDELER_GRAPHQL_API_KEY
DATAFORDELEREN_API_KEY
```

## GSearch request template

```ts
const url = new URL("https://api.dataforsyningen.dk/rest/gsearch/v2.0/husnummer");
url.searchParams.set("token", process.env.GSEARCH_TOKEN ?? "");
url.searchParams.set("q", query);
url.searchParams.set("limit", "10");
url.searchParams.set("srid", "4326");
```

Never hardcode a token in source code. Use an environment variable, secret
manager, or server-side proxy.

## Resources

Use singular GSearch resource names:

```txt
adresse
husnummer
navngivenvej
stednavn
kommune
region
retskreds
postnummer
opstillingskreds
sogn
politikreds
matrikel
matrikel_udgaaet
```

## Migration heuristics

- DAWA `/adresser/autocomplete` becomes GSearch `adresse`.
- DAWA `/adgangsadresser/autocomplete` becomes GSearch `husnummer`.
- DAWA road autocomplete becomes GSearch `navngivenvej`.
- DAWA administrative autocomplete endpoints usually become the singular GSearch
  resource name: `kommuner` to `kommune`, `regioner` to `region`, and so on.
- DAWA `reverse` does not become a direct GSearch endpoint. For nearest-address
  behavior, use GSearch ECQL spatial filters (`INTERSECTS`, `DWITHIN`) with
  EPSG:25832 coordinates; for authoritative reverse semantics, use a documented
  service or your own spatial index built from authoritative data.
- DAWA `datavask`, `replikering`, history, BBR, BFE, and bulk data do not become
  GSearch. Pick DAR/Datafordeler, the planned Adressevask replacement when it is
  available, events, downloads, or another authoritative service.
- DAWA `bygninger` does not become GSearch. For targeted building geometry, use
  GeoDanmark Vektor GraphQL `GEODKV_Bygning` through Datafordeler with
  EPSG:25832 WKT spatial filters.
- DAWA ID lookup endpoints should not be replaced with guessed GSearch path
  lookups. Use DAR/Datafordeler or documented Adressevælger lookup paths.

## Code review checklist

Flag these issues:

- new code uses `dawa.aws.dk`;
- new code uses `api.dataforsyningen.dk/adresser` or `/adgangsadresser` for
  autocomplete instead of GSearch or Adressevælger;
- token is hardcoded;
- code parses `visningstekst` instead of using structured fields;
- code stores only a display label and loses the UUID;
- code selects `adresse` results but fails to preserve or derive the parent
  `husnummer` UUID for BBR/BFE workflows;
- `adresse` and `husnummer` are confused;
- spatial filters use WGS84 coordinates instead of EPSG:25832;
- examples call GSearch "Google Search".

## Suggested agent answer

When asked "use the Danish address API", answer along these lines:

> Use SDFI/Dataforsyningen GSearch for Danish address search. Query
> `/rest/gsearch/v2.0/husnummer` for building-level suggestions and
> `/rest/gsearch/v2.0/adresse` for unit-level suggestions. Store the returned
> DAR UUID and structured fields; do not parse the display string. Preserve the
> parent `husnummer` UUID when a unit-level `adresse` is selected. Use
> DAR/Datafordeler/Adressevælger and the planned Adressevask replacement for
> non-search DAWA features.
