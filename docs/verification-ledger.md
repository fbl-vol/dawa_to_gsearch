# Verification ledger

This ledger records the main factual claims in this repo, the source used to
verify each one, and any live checks run from this workspace.

Reviewed on **2026-06-27**.

Confluence pages are official sources but can be flaky from command-line curl
checks. The Adressevælger/Adressevask claims below were verified from the
official Klimadatastyrelsen/SDFI pages during review; if updating them later,
re-open the pages in a browser and re-check the status text.

## Source-backed claims

| Claim | Verification | Status |
| --- | --- | --- |
| GSearch is a REST API for searching addresses, cadastral parcels, DAGI geography, and Danish place names. | Official GSearch docs list the resource families and describe typeahead/search usage. | Verified |
| GSearch uses `https://api.dataforsyningen.dk/rest/gsearch/v2.0/{resource}`. | Official GSearch docs show this request syntax. | Verified |
| `q` is the central search parameter; `limit` defaults to 10 and maxes at 100. | Official GSearch docs describe `q` and `limit`. | Verified |
| GSearch filters use ECQL. | Official GSearch docs refer to ECQL and GeoServer CQL/ECQL syntax. | Verified |
| Spatial filters use EPSG:25832. | Official GSearch docs state that geometries in filters must use EPSG:25832. | Verified |
| Response objects include `id`, `visningstekst`, and GeoJSON geometry. | Official GSearch docs describe response contents and geometry. | Verified |
| GSearch resources include `adresse`, `husnummer`, `kommune`, `matrikel`, `matrikel_udgaaet`, `navngivenvej`, `opstillingskreds`, `politikreds`, `postnummer`, `region`, `retskreds`, `sogn`, and `stednavn`. | Official GSearch docs and GSearch UI default resource list. | Verified |
| Adressevælger replaces the value proposition of DAWA autocomplete, but is not a 1:1 interface/response replacement. | Official Adressevælger FAQ. | Verified |
| Adressevælger does not provide DAWA-style reverse geocoding. | Official Adressevælger FAQ points reverse geocoding users to Datafordeler options. | Verified |
| Adressevælger documents ID lookup paths for house numbers and addresses. | Official Adressevælger ID lookup page. | Verified |
| Adressevask is the planned replacement for DAWA address `datavask`. | Official Adressevælger & Adressevask page. | Verified |
| Adressevask was not available on the reviewed page and was expected at the end of August 2026. | Official Adressevælger & Adressevask page. | Verified |
| DAWA BBR data has been phased out from DAWA. | Official DAWA BBR documentation. | Verified |
| DAR BFE can fetch BFE numbers from an Adresse or Husnummer. | Datafordeler DAR BFE page. | Verified |
| DAR BFE REST is marked for phase-out by Datafordeler. | Datafordeler DAR BFE page marks REST as being phased out. | Verified |
| DAR GraphQL exposes address entities and uses API-key/OAuth access. | Datafordeler DAR GraphQL page. | Verified |
| BBR GraphQL exposes BBR entities and uses API-key/OAuth access. | Datafordeler BBR GraphQL page. | Verified |
| MAT GraphQL exposes Matriklen entities and uses API-key/OAuth access. | Datafordeler MAT GraphQL page. | Verified |
| BBR REST is marked for phase-out by Datafordeler. | Datafordeler BBR REST overview. | Verified |
| DAR and BBR events are GraphQL-accessible event entities. | Datafordeler DAR Hændelser and BBR Hændelser pages. | Verified |
| DAR file download provides pre-generated entity files in JSON/CSV with API-key/OAuth access. | Datafordeler DAR Fildownload page. | Verified |
| DAR WFS entities provide GeoJSON/GML access and expose current/historical layers. | Datafordeler DAR WFS entiteter page. | Verified |
| Datafordeler GeoDanmark Vektor GraphQL exposes entity data including `Bygning` and uses API-key/OAuth access. | Datafordeler GeoDanmark Vektor GraphQL page. | Verified |
| Datafordeler GeoDanmark Vektor GraphQL schema can be used to inspect entity fields. | Datafordeler GeoDanmark Vektor GraphQL schema page and live schema probe. | Verified |
| `GEODKV_Bygning` includes building geometry and BBR/DAR-adjacent identifiers useful for migration: `geometri`, `bygningstype`, `status`, `id_lokalId`, and `BBRUUID`. | Live `https://graphql.datafordeler.dk/GEODKV/v2/schema` probe with local API key. | Verified |
| `GEODKV_BygningFilterInput.geometri.intersects` accepts EPSG:25832 WKT via `SpatialIntersectsOperationFilterInput`. | Live `GEODKV` GraphQL schema probe with local API key. | Verified |
| GeoDanmark Vektor Fildownload includes the `Bygning` entity for pre-generated/bulk workflows. | Datafordeler GeoDanmark Vektor Fildownload page. | Verified |
| GeoDanmark Vektor WFS entities are available as generated entity layers with GeoJSON/GML output. | Datafordeler GeoDanmark Vektor WFS entities page. | Verified |
| The older GeoDanmark Vektor WFS page reviewed on 2026-06-27 marks that specific WFS service for phase-out in 2026, and old GeoDanmark `filudtræk` documentation points users toward Fildownload. | Datafordeler GeoDanmark Vektor legacy WFS page and GeoDanmark `filudtræk` page. | Verified |
| Dataforsyningen webservices require HTTPS and generally use `token=<token>` for authenticated calls. | Dataforsyningen webservice guide. | Verified |

## Live checks

These checks used credentials from the local private `dk-building-data` repo.
No secret values were printed or copied into this repository.

| Check | Result |
| --- | --- |
| GSearch `husnummer` with `q=genvej` | HTTP 200, first result `Genvej 1, 4874 Gedser` |
| GSearch `adresse` with `q=flens` | HTTP 200, first result `Flensbjerg 1A, 4960 Holeby` |
| GSearch `navngivenvej` with `q=krin` | HTTP 200, first result `Kringelborg Alle (4800 Nykøbing F)` |
| GSearch `postnummer` with `q=mari` | HTTP 200, first result `9550 Mariager` |
| GSearch `kommune` with `q=aalborg` | HTTP 200, first result `Aalborg Kommune` |
| GSearch `region` with `q=midt` | HTTP 200, first result `Region Midtjylland` |
| GSearch `sogn` with `q=budolfi` | HTTP 200, first result `Budolfi sogn` |
| GSearch `politikreds` with `q=vest` | HTTP 200, first result `Københavns Vestegns Politikreds` |
| GSearch `retskreds` with `q=københavn` | HTTP 200, first result `Københavns Byret` |
| GSearch `opstillingskreds` with `q=vest` | HTTP 200, first result `Vesterbrokredsen` |
| GSearch `matrikel` with `q=123ab` | HTTP 200, first result `123ab, Povlsker` |
| GSearch `matrikel_udgaaet` with `q=11a` | HTTP 200, first result `11ae, Em By, Em` |
| GSearch `stednavn` with `q=Benedikte` | HTTP 200, first result `Benedikte Sø (Sø i Gråsten)` |
| GSearch `postnummer` spatial filter `INTERSECTS(geometri,POINT(689255 6051787))` | HTTP 200, returned `4874 Gedser` |
| GSearch `husnummer` spatial filter `DWITHIN(geometri,POINT(689255 6051787),100,meters)` | HTTP 200, returned nearby `husnummer` candidates |
| GSearch `/adresse/{id}` path lookup | HTTP 404 in live check; repo guidance avoids treating it as DAWA parity |
| GSearch `/husnummer/{id}` path lookup | HTTP 404 in live check; repo guidance avoids treating it as DAWA parity |
| DAR BFE dummy `HusnummerId` request | HTTP 200 with empty `jordstykkeList`, confirming endpoint shape and authentication path |
| Datafordeler DAR GraphQL `query { __typename }` | HTTP 200 with `{"data":{"__typename":"Query"}}` |
| Datafordeler BBR GraphQL `query { __typename }` | HTTP 200 in executable example verification |
| Datafordeler MAT GraphQL `query { __typename }` | HTTP 200 in executable example verification |
| Datafordeler GEODKV GraphQL schema with local API key | HTTP 200; schema contains `GEODKV_Bygning`, `geometri`, `BBRUUID`, `id_lokalId`, and spatial filter inputs |
| Datafordeler GEODKV `GEODKV_Bygning` spatial query near Gedser | HTTP 200; returned 3 building candidates, first with `bygningstype=Bygning`, `status=Anlagt`, `crs=25832` |
| Python GSearch address helper | Returned 3 merged suggestions for `Søbakkevej 8, Tilst`; first result had `id`, `husnummerId`, and numeric longitude/latitude |
| Python GEODKV building-footprint helper | Returned 3 `GEODKV_Bygning` candidates near Gedser; first had WKT geometry, `status=Anlagt`, `crs=25832`, `id_lokalId`, and `BBRUUID` |
| Go GSearch address helper | Returned 3 merged suggestions for `Søbakkevej 8, Tilst`; first result had `id`, `husnummerId`, and numeric longitude/latitude |
| Go GEODKV building-footprint helper | Returned 3 `GEODKV_Bygning` candidates near Gedser; first had WKT geometry, `status=Anlagt`, `crs=25832`, `id_lokalId`, and `BBRUUID` |

## Endpoint coverage

The generated DAWA operation list was extracted from
<https://docs.dataforsyningen.dk/> and compared mechanically with
[DAWA endpoint coverage](dawa-endpoint-coverage.md).

Result on 2026-06-27:

```txt
remote_count=191
doc_count=191
missing_from_doc=0
extra_in_doc=0
```

## Deliberate non-claims

- This repo does not claim GSearch replaces every DAWA endpoint.
- This repo does not claim GSearch provides DAWA-style reverse geocoding.
- This repo does not claim GSearch is a BBR, BFE, history, event, replication,
  or local-copy service.
- This repo does not claim Adressevask is available before the official
  Klimadatastyrelsen/SDFI page says it is.
- This repo does not publish authenticated GSearch or Datafordeler responses
  that contain secret material.
