# DAWA endpoint target mapping

This mapping is intentionally job-based. If a DAWA feature is not clearly a
search/typeahead feature, migrate it to DAR/Datafordeler, Adressevælger,
Adressevask, events, downloads, WFS/OGC, or a dataset-specific API instead of
forcing it through GSearch.

For an audit against the generated DAWA operation list, see
[DAWA endpoint coverage](dawa-endpoint-coverage.md).

## Address resources

| DAWA endpoint or concept | GSearch resource | Migration note |
| --- | --- | --- |
| `/adresser/autocomplete` | `adresse` | Unit-level addresses. Includes floor/door fields where present. |
| `/adresser?q=...` | `adresse` for search; DAR/Datafordeler for full register querying | GSearch is for best-match search, not arbitrary DAWA query parity. |
| `/adresser/{id}` | DAR/Datafordeler or Adressevælger ID lookup | GSearch is the search step. Live GSearch path lookups returned 404 during verification, so do not rely on `/adresse/{id}` as DAWA parity. |
| `/adgangsadresser/autocomplete` | `husnummer` | Building-level/access-address suggestions. |
| `/adgangsadresser?q=...` | `husnummer` for search; DAR/Datafordeler for full register querying | Keep `husnummer` UUID for downstream building lookups. |
| `/adgangsadresser/{id}` | DAR/Datafordeler or Adressevælger ID lookup | GSearch is the search step. Live GSearch path lookups returned 404 during verification, so do not rely on `/husnummer/{id}` as DAWA parity. |

## Roads and postcodes

| DAWA endpoint or concept | GSearch resource | Migration note |
| --- | --- | --- |
| `/vejnavne/autocomplete` | `navngivenvej` | Named-road search. |
| `/vejstykker/autocomplete` | `navngivenvej` | GSearch groups this under named roads. |
| `/navngivneveje/autocomplete` | `navngivenvej` | Singular resource name. |
| `/postnumre/autocomplete` | `postnummer` | Postcode district search. |
| `/supplerendebynavne*` | No direct GSearch default resource | Check returned address fields or use an authoritative dataset. |

## DAGI administrative geography

| DAWA endpoint or concept | GSearch resource |
| --- | --- |
| `/kommuner/autocomplete` | `kommune` |
| `/regioner/autocomplete` | `region` |
| `/sogne/autocomplete` | `sogn` |
| `/politikredse/autocomplete` | `politikreds` |
| `/retskredse/autocomplete` | `retskreds` |
| `/opstillingskredse/autocomplete` | `opstillingskreds` |
| `/afstemningsomraader*` | No direct GSearch default resource |
| `/storkredse*` | No direct GSearch default resource |
| `/valglandsdele*` | No direct GSearch default resource |
| `/landsdele*` | No direct GSearch default resource |

## Cadastral and place-name data

| DAWA endpoint or concept | GSearch resource | Migration note |
| --- | --- | --- |
| `/jordstykker/autocomplete` | `matrikel` | Current cadastral parcel search. |
| Retired cadastral parcel search | `matrikel_udgaaet` | Use only when retired parcels are expected. |
| `/ejerlav/autocomplete` | Partial via `matrikel` filters | GSearch does not expose a standalone default `ejerlav` resource in the main resource list. |
| Place-name search | `stednavn` | Supports type/subtype filters such as `stednavn_type` where relevant. |

## Features that use non-GSearch targets

| DAWA feature | Use instead |
| --- | --- |
| DAWA autocomplete UI successor | Adressevælger or your own UI on top of GSearch |
| Reverse geocoding | No direct GSearch reverse endpoint; use a documented reverse/spatial service or custom lookup from authoritative data. GSearch ECQL spatial filters can support nearest-address lookup. |
| Address washing/data normalization | Planned Adressevask replacement when available, or app-specific validation |
| Replication/change feeds | Datafordeler events/download services |
| Historical queries | DAR/Datafordeler or dataset-specific historical services |
| BBR/BFE | GSearch `husnummer` selection, then DAR/Datafordeler BFE/BBR/MAT services |
| Building polygons/footprints | Datafordeler GeoDanmark Vektor GraphQL `GEODKV_Bygning` for targeted lookup; GeoDanmark Vektor Fildownload/WFS entities for bulk/GIS |
| Bulk downloads | Dataforsyningen download services or Datafordeler |

## GeoSearch to GSearch note

The upstream GSearch repo also documents migration from SDFI's older GeoSearch.
Those resource names map as follows:

| GeoSearch | GSearch |
| --- | --- |
| `adresser` | `adresse` |
| `kommuner` | `kommune` |
| `matrikelnumre` | `matrikel` |
| `matrikelnumre_udgaaet` | `matrikel_udgaaet` |
| `opstillingskredse` | `opstillingskreds` |
| `politikredse` | `politikreds` |
| `postdistriker` | `postnummer` |
| `regioner` | `region` |
| `retskredse` | `retskreds` |
| `sogne` | `sogn` |
| `stednavn_v3` | `stednavn` |
| no GeoSearch equivalent | `husnummer`, `navngivenvej` |
