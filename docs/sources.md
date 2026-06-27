# Verified sources

Sources reviewed on 2026-06-27.

Most non-Confluence links below passed direct `curl -L` verification.
Confluence pages are official SDFI/Klimadatastyrelsen sources, but they can be
flaky from command-line checks. They opened successfully in a browser-style web
check on 2026-06-27. Re-open them manually in a browser if you update claims
from them.

## Official sources

| Source | URL | Why it matters |
| --- | --- | --- |
| DAWA documentation | <https://dawadocs.dataforsyningen.dk/> | States that DAWA is closing and describes DAWA's scope: addresses, DAGI, cadastral data, BBR, and Danish place names. |
| DAWA API documentation | <https://dawadocs.dataforsyningen.dk/dok/api> | Lists DAWA endpoints, including autocomplete, reverse, datavask, replication, and ID lookups. |
| DAWA BBR documentation | <https://dawadocs.dataforsyningen.dk/dok/bbr> | States that DAWA BBR data is no longer updated from 2023-12-16 and was fully phased out from DAWA on 2024-04-01. |
| DAWA shutdown page | <https://dataforsyningen.dk/data/4924> | Official shutdown information. The site is JavaScript-rendered, so verify manually in a browser for the latest date. |
| GSearch repository | <https://github.com/Klimadatastyrelsen/gsearch> | Upstream implementation and README. Documents GSearch principles and GeoSearch-to-GSearch mapping. |
| GSearch documentation | <https://github.com/Klimadatastyrelsen/gsearch/tree/v2.0/doc> | Documents resources, request syntax, `q`, `limit`, ECQL `filter`, response geometry, and examples. |
| GSearch UI repository | <https://github.com/Klimadatastyrelsen/gsearch-ui> | Official web component. Documents required `data-token`, default resources, `data-limit`, `data-api`, `data-filter`, and `data-srid`. |
| Dataforsyningen webservice guide | <https://confluence.sdfi.dk/pages/viewpage.action?pageId=158368397> | Documents HTTPS, `api.dataforsyningen.dk`, `/rest/` prefix, and `token=<token>` usage. |
| Adressevælger & Adressevask | <https://confluence.sdfi.dk/pages/viewpage.action?pageId=234782998> | Klimadatastyrelsen documentation for Adressevælger and Adressevask. States Adressevælger replaces DAWA autocomplete and that Adressevask is the planned replacement for DAWA's address datavask service, expected at the end of August 2026 and not yet available as of the reviewed page. |
| Adressevælger FAQ | <https://confluence.sdfi.dk/display/ADV/FAQ> | States that Adressevælger replaces DAWA autocomplete value, differs from DAWA's response/interface, returns DAR UUIDs, has API and JavaScript component, has no DAWA-style reverse geocoding, and requires a token parameter. |
| Adressevælger ID lookup | <https://confluence.sdfi.dk/pages/viewpage.action?pageId=246743156> | Documents Adressevælger id lookup for `husnumre/{id}` and `adresser/{id}`, and says husnummer/adresse IDs can be used for GraphQL searches in Datafordeler when other register data is needed. |
| GeoServer ECQL guide | <https://docs.geoserver.org/main/en/user/tutorials/cql/cql_tutorial/> | Explains CQL/ECQL syntax used by GSearch filters. |
| Dataforsyningen OpenAPI/docs portal | <https://docs.dataforsyningen.dk/> | Contains generated API and schema documentation, including DAWA and GSearch schema sections. |
| Datafordeler DAR BFE | <https://datafordeler.dk/dataoversigt/danmarks-adresseregister-dar/dar-bfe/> | Documents the DAR BFE service for fetching BFE numbers from an Adresse or Husnummer, including the REST service URL. |
| Datafordeler DAR GraphQL | <https://datafordeler.dk/dataoversigt/danmarks-adresseregister-dar/dar-graphql/> | Documents DAR entities such as Adresse, Husnummer, NavngivenVej, Postnummer, and SupplerendeBynavn through Datafordeler GraphQL. |
| Datafordeler BBR GraphQL | <https://datafordeler.dk/dataoversigt/bygnings-og-boligregistret-bbr/bbr-graphql/> | Documents BBR entities such as Bygning, Enhed, Etage, Grund, Opgang, and TekniskAnlæg through Datafordeler GraphQL. |
| Datafordeler BBR REST overview | <https://datafordeler.dk/dataoversigt/bygnings-og-boligregistret-bbr/bbr/> | Documents the legacy BBR REST service and BBR methods such as bygning, enhed, grund, and tekniskanlaeg. |
| Datafordeler MAT GraphQL | <https://datafordeler.dk/dataoversigt/matriklen-mat/matriklen-graphql/> | Documents Matriklen/MAT GraphQL entities and API-key/OAuth access. |
| Datafordeler DAR events | <https://datafordeler.dk/dataoversigt/danmarks-adresseregister-dar/dar-haendelser/> | Documents `DAR_Events` as the Datafordeler event entity for DAR. |
| Datafordeler BBR events | <https://datafordeler.dk/dataoversigt/bygnings-og-boligregistret-bbr/bbr-haendelser/> | Documents `BBR_Events` as the Datafordeler event entity for BBR. |
| Datafordeler DAR file download | <https://datafordeler.dk/dataoversigt/danmarks-adresseregister-dar/dar-fildownload/> | Documents pre-generated DAR entity downloads in JSON/CSV with API-key/OAuth access. |
| Datafordeler DAR WFS entities | <https://datafordeler.dk/dataoversigt/danmarks-adresseregister-dar/dar-wfs-entiteter/> | Documents DAR entity WFS access, GeoJSON/GML output, GetCapabilities, and API-key/OAuth access. |
| GeoDanmark Vektor GraphQL | <https://datafordeler.dk/dataoversigt/geodanmark-vektor/geodanmark-vektor-graphql/> | Documents GeoDanmark Vektor as a GraphQL Query service with entity data including `Bygning`, API-key/OAuth access, and Datafordeler Administration setup. |
| GeoDanmark Vektor GraphQL schema | <https://datafordeler.dk/dataoversigt/geodanmark-vektor/geodanmark-vektor-graphql-skema/> | Documents the GeoDanmark Vektor GraphQL schema service. Live schema probing verified `GEODKV_Bygning`, `geometri`, `BBRUUID`, `id_lokalId`, and spatial filter inputs. |
| GeoDanmark Vektor WFS entities | <https://datafordeler.dk/dataoversigt/geodanmark-vektor/geodanmark-vektor-wfs-entiteter/> | Documents generated entity-based WFS layers for GeoDanmark Vektor with GeoJSON/GML output and API-key/OAuth access. |
| GeoDanmark Vektor Fildownload | <https://datafordeler.dk/dataoversigt/geodanmark-vektor/geodanmark-vektor-fildownload/> | Documents pre-generated GeoDanmark Vektor entity files including `Bygning`, with API-key/OAuth access. |
| GeoDanmark Vektor legacy WFS | <https://datafordeler.dk/dataoversigt/geodanmark-vektor/geodanmark-vektor-wfs/> | Documents older GeoDanmark vector WFS, including building categories; the reviewed page marks that specific WFS service for phase-out in 2026. |
| GeoDanmark legacy filudtræk | <https://confluence.sdfi.dk/x/cYWIAQ> | States that old GeoDanmark `filudtræk` is being phased out and points users toward Fildownload. |

## Implementation notes

Useful implementation details from existing GSearch client code reviewed during
setup:

- `GSEARCH_TOKEN` is required.
- The base URL is `https://api.dataforsyningen.dk/rest/gsearch/v2.0`.
- Address search queries both `husnummer` and `adresse`.
- `srid=4326` is used when longitude/latitude output is needed.
- Results are deduplicated while preserving the parent `husnummer` UUID.
- Address-led BFE/BBR lookup can use GSearch `husnummer` UUIDs as the bridge
  into DAR BFE/Datafordeler, then continue through BBR/MAT/DAR services.
- Nearest-address reverse lookup can be implemented with GSearch ECQL spatial
  filters, but this is a custom lookup pattern, not a DAWA `/reverse` clone.
- Targeted building-footprint lookup should use Datafordeler GeoDanmark Vektor
  GraphQL `GEODKV_Bygning` with an EPSG:25832 WKT spatial filter; bulk/local
  copies should use GeoDanmark Vektor Fildownload or WFS entities.

No secrets or private client code are copied here.

## Live verification notes

Unauthenticated GSearch requests currently return:

```json
{"message":"Authentication token must only be lowercase (token) or uppercase (TOKEN)!","http_status_code":400}
```

Using an authorized local `GSEARCH_TOKEN` as a query parameter returned HTTP 200
for all documented GSearch resources listed in
[examples gallery](examples-gallery.md). Direct GSearch path lookups tested as
`/adresse/{id}` and `/husnummer/{id}` returned HTTP 404, so this repo does not
recommend those path forms as DAWA `{id}` parity.

Using a lowercase `token` header in earlier direct curl tests returned the
gateway's `400` token-parameter error. Because of that, this repo's public
examples use the documented `token` query parameter and describe the header form
only as an environment-specific pattern to verify.
