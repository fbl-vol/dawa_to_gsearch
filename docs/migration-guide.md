# DAWA migration guide

The safest migration path is to split each DAWA call by job and move it to the
right post-DAWA service. GSearch is the search/typeahead path. DAR/Datafordeler,
Dataforsyningen downloads, events, WFS/OGC, Adressevælger, Adressevask, and
dataset-specific APIs cover the rest.

Use GSearch where the job is search. Use a register/data service where the job
is authoritative lookup, history, replication, bulk data, geometry, or BBR/BFE.

## Migration workflow

1. Inventory every DAWA URL your app calls.
2. Classify each call as autocomplete/search, lookup by ID, reverse geocoding,
   data washing, replication/history, BBR, or bulk data.
3. Move autocomplete/search calls to GSearch or Adressevælger.
4. Move DAWA autocomplete UI needs to either your own GSearch UI or SDFI's
   Adressevælger component/API.
5. Move authoritative register, BFE/BBR, and historical workflows to
   DAR/Datafordeler or dataset-specific services.
6. Add contract tests using saved DAWA examples and expected GSearch outputs.
7. Deploy behind a feature flag and log no-result/auth/error rates.

## Search paths

| DAWA usage | Search replacement | Notes |
| --- | --- | --- |
| `/adresser/autocomplete` | `/rest/gsearch/v2.0/adresse` | Use for unit-level address selection. |
| `/adgangsadresser/autocomplete` | `/rest/gsearch/v2.0/husnummer` | Use for building/access-address selection. |
| `/vejnavne/autocomplete`, `/vejstykker/autocomplete`, `/navngivneveje/autocomplete` | `/rest/gsearch/v2.0/navngivenvej` | GSearch's resource name is singular. |
| `/postnumre/autocomplete` | `/rest/gsearch/v2.0/postnummer` | Returns postcode district matches. |
| `/kommuner/autocomplete` | `/rest/gsearch/v2.0/kommune` | DAGI municipality search. |
| `/regioner/autocomplete` | `/rest/gsearch/v2.0/region` | DAGI region search. |
| `/sogne/autocomplete` | `/rest/gsearch/v2.0/sogn` | DAGI parish search. |
| `/politikredse/autocomplete` | `/rest/gsearch/v2.0/politikreds` | DAGI police district search. |
| `/retskredse/autocomplete` | `/rest/gsearch/v2.0/retskreds` | DAGI court district search. |
| `/opstillingskredse/autocomplete` | `/rest/gsearch/v2.0/opstillingskreds` | DAGI election district search. |
| `/jordstykker/autocomplete`, cadastral search | `/rest/gsearch/v2.0/matrikel` or `/matrikel_udgaaet` | Use current or retired cadastral resources as appropriate. |
| DAWA place-name search | `/rest/gsearch/v2.0/stednavn` | Search Danish place names. |

## Register and data paths

| DAWA feature | Recommended replacement |
| --- | --- |
| Reverse geocoding endpoints | No direct GSearch reverse endpoint. Use a documented reverse/spatial service or build a custom spatial lookup from authoritative data. For nearest-address behavior, GSearch can be used with EPSG:25832 ECQL filters such as `INTERSECTS` and `DWITHIN`; see [implementation patterns](implementation-patterns.md). |
| `datavask` address normalization | Klimadatastyrelsen/SDFI's planned Adressevask replacement, or your own normalization and validation workflow. As of 2026-06-27, the official docs say Adressevask is not available yet and is expected at the end of August 2026. |
| Replication endpoints | Datafordeler events, downloads, WFS, or dataset-specific APIs. |
| History endpoints | DAR/Datafordeler or authoritative historical datasets where available. |
| BBR through DAWA | Use GSearch to select a `husnummer` UUID, then DAR/Datafordeler BFE/BBR/MAT services. DAWA BBR data has already been phased out. |
| Building polygons through DAWA | Datafordeler GeoDanmark Vektor GraphQL `GEODKV_Bygning` for targeted lookup; GeoDanmark Vektor Fildownload/WFS entities for bulk or GIS workflows. |
| Bulk exports/downloads | Dataforsyningen download services or Datafordeler, depending on dataset. |

## Architecture pattern

Most migrations should become a two-step flow:

1. Let the user find the thing with GSearch or Adressevælger.
2. Use the selected IDs to fetch authoritative data from DAR/Datafordeler,
   BBR/MAT, events, downloads, WFS/OGC, or the dataset-specific API.

This keeps search, selection, authoritative lookup, and local-copy workflows
separate. It also keeps code agents from forcing non-search DAWA features
through GSearch.

## Data model changes to expect

DAWA and GSearch both expose DAR UUIDs, but response shapes and field names are
not identical. Update your application model explicitly.

Recommended selected-address model:

```ts
type SelectedDanishAddress = {
  source: "gsearch";
  resource: "adresse" | "husnummer";
  id: string;
  husnummerId?: string;
  label: string;
  kommunekode?: string;
  vejkode?: string;
  postnummer?: string;
  longitude?: number;
  latitude?: number;
};
```

For an autocomplete that queries both `adresse` and `husnummer`, preserve the
`husnummer` UUID even when the selected display result is a unit-level
`adresse`. Building, parcel, and BBR workflows often need the building/access
address identifier rather than the unit identifier.

For BFE/BBR workflows, treat the GSearch result as the address-selection step.
Resolve the `husnummer` UUID through DAR/Datafordeler before querying BBR,
property, cadastral, or building-footprint data.

## Query changes

DAWA examples often use endpoint-specific parameters. GSearch mostly uses:

- `q` for the search string;
- `limit` for result count;
- `filter` for ECQL constraints;
- `srid` for output geometry coordinate system;
- `token` for authentication.

Replace DAWA query composition with a small GSearch client instead of scattering
URL strings through the codebase.

## UI migration pattern

1. Render suggestions from `visningstekst`.
2. Keep the full result object in memory while the suggestion list is open.
3. On selection, store `id`, `resource`, `visningstekst`, coordinates, and the
   structured fields your downstream systems need.
4. Do not store only the display string.
5. If you need a hosted component, evaluate `@dataforsyningen/gsearch-ui` or
   SDFI Adressevælger.

## Testing checklist

Use real examples from your application:

- ordinary street address;
- address with floor/door;
- address with Danish characters such as `å`, `ø`, and `æ`;
- old spelling or phonetic variant;
- postcode-only search;
- road search without house number;
- municipality-filtered result;
- no-result case;
- invalid/expired token;
- broad query with `limit=100`;
- spatial filter if your app uses maps.

For each example, assert the selected result's UUID and structured fields, not
just the display label.

## Deployment checklist

- Token is configured through environment/secret management.
- No DAWA URLs remain in new search/autocomplete code.
- Existing DAWA URLs are listed in a temporary migration register with owner and
  removal date.
- Errors distinguish auth (`401`), bad request/filter (`400`), no results, and
  timeout (`504`).
- Metrics track query count, no-result rate, latency, auth failures, and selected
  resource type.
- Documentation tells developers that GSearch is SDFI GSearch, not Google
  Search.
