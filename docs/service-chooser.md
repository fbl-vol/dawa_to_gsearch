# Service chooser

DAWA bundled search, address lookup, register-ish convenience endpoints,
replication, history, data washing, and BBR/OIS access behind one familiar API.
The replacement landscape is more explicit. Pick the service by job.

If you only remember one rule: use GSearch for user-facing search, then switch
to the authoritative service for register data, history, replication, geometry,
or normalization.

## Quick decision table

| Job | Use | Why |
| --- | --- | --- |
| Build your own Danish address autocomplete | GSearch `husnummer` + `adresse` | Direct search API, flexible UI, supports filters and multiple resources. |
| Use an official hosted address selector | Adressevælger | SDFI's address picker/API path for the DAWA autocomplete value proposition. |
| Search roads, postcodes, municipalities, regions, parishes, cadastral parcels, or place names | GSearch | These are first-class GSearch resources. |
| Validate/normalize messy address text | Adressevask when available, or a dedicated validation flow | GSearch is search, not data washing. The official Adressevask replacement is expected at the end of August 2026 and is not available yet as of 2026-06-27. |
| Fetch full DAR/register objects | DAR/Datafordeler GraphQL, WFS, downloads, or relevant authoritative service | GSearch returns search results, not every authoritative register view. |
| Resolve an address to BFE/BBR/property data | GSearch, then DAR/Datafordeler | Use GSearch for `husnummer` selection, then resolve through DAR BFE/BBR/MAT. |
| Fetch building footprints | Datafordeler GeoDanmark Vektor GraphQL for targeted lookup; GeoDanmark Vektor Fildownload/WFS entities for bulk/GIS | GSearch is not a building-geometry source. |
| Reverse geocode coordinates to an address | Documented reverse/spatial service, or custom spatial lookup from authoritative data | Adressevælger FAQ says DAWA-style reverse geocoding is not provided there; GSearch is not a direct reverse endpoint, though ECQL spatial filters can support nearest-address workflows. |
| Subscribe to changes, replicate, or build a local copy | Datafordeler events, file downloads, WFS, OGC services, or dataset exports | GSearch is an online search API, not a replication or local-copy feed. |
| Search the public web | Google Search or another web search API | GSearch is not Google Search. |

## When GSearch is the right default

Choose GSearch when the user types text and your app needs likely official
Danish geodata matches.

Good examples:

- checkout address autocomplete;
- CRM address picker;
- municipality-filtered address search;
- road-name search;
- postcode or municipality lookup;
- map search for place names or cadastral parcels.

Implementation pattern:

1. Query the relevant resource with `q`.
2. Use `limit` to cap suggestions.
3. Add `filter` only after the unfiltered search works.
4. Store `id`, `resource`, `visningstekst`, and structured fields.

## When Adressevælger is better

Choose Adressevælger when you want SDFI's official address selector experience
instead of owning the search UI yourself.

Expect changes from DAWA autocomplete. The official FAQ says Adressevælger
returns DAR UUIDs like DAWA autocomplete did, but the response/interface is not
a 1:1 replacement and local adaptations should be expected.

## What Adressevask means

Adressevask is Klimadatastyrelsen/SDFI's planned replacement for DAWA's
address `datavask` service. It is meant for messy address input: it can take a
text address or dedicated address fields and return the best matching official
addresses. The official documentation currently says the new service is not
available yet and is expected at the end of August 2026.

## When DAR/Datafordeler is better

Choose DAR/Datafordeler when the app needs authoritative register data after a
selection has been made.

Typical flow:

1. Use GSearch or Adressevælger for user selection.
2. Persist the DAR UUID and structured fields.
3. Use DAR/Datafordeler/BBR/MAT services for downstream register lookups.

This keeps the UI search problem separate from the authoritative data retrieval
problem.

See [implementation patterns](implementation-patterns.md) for the common
GSearch `husnummer` UUID to DAR BFE/BBR flow and the custom reverse-geocode
pattern.

See [DAR and Datafordeler guide](datafordeler-dar-guide.md) for the
authoritative-data side: BFE, BBR/MAT, GraphQL, events, downloads, WFS/OGC, and
local-copy replacement paths.

## When not to use GSearch

Do not use GSearch as a substitute for:

- DAWA `datavask`;
- DAWA `reverse`;
- DAWA replication feeds;
- BBR/BFE data access;
- bulk exports;
- arbitrary full-register querying.

Those workflows need a different service even if the old implementation used
DAWA.
