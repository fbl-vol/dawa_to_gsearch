# DAWA endpoint coverage

This page checks the generated DAWA API operation list, not just the common
autocomplete endpoints. It is meant to answer: "Does every DAWA endpoint family
have a migration path or an explicit explanation?"

Coverage source:

- Human DAWA API overview: <https://dawadocs.dataforsyningen.dk/dok/api>
- Generated API docs portal: <https://docs.dataforsyningen.dk/>
- Extracted generated DAWA operation count on 2026-06-27: **191 GET operations**

## Migration classes

| Class | Meaning |
| --- | --- |
| GSearch direct | Use GSearch for the same search/autocomplete job. |
| GSearch partial | GSearch helps with search, but lookup/reverse/history/full data continues through DAR/Datafordeler, downloads, WFS/OGC, or another authoritative path. |
| Custom GSearch spatial lookup | Use GSearch ECQL spatial filters for a nearest-object workflow, without claiming DAWA endpoint parity. |
| Adressevælger/Adressevask | Use SDFI/Klimadatastyrelsen's address picker or planned address washing path instead of GSearch. |
| DAR/Datafordeler/authoritative service | Use DAR/Datafordeler, OGC/WFS/download services, or a dataset-specific API. |
| No GSearch equivalent | Do not migrate this endpoint to GSearch; document a separate replacement. |

## Address and address-selection endpoints

### `adgangsadresser`

Generated operations:

```txt
adgangsadresser
adgangsadresser_autocomplete
adgangsadresser_reverse
adgangsadresser_{id}
historik_adgangsadresser
datavask_adgangsadresser
```

Migration path:

- `adgangsadresser_autocomplete`: **GSearch direct** via `husnummer`.
- `adgangsadresser` text search: **GSearch partial** via `husnummer` for
  best-match search; use DAR/Datafordeler for full query semantics.
- `adgangsadresser_{id}`: **DAR/Datafordeler/authoritative service** for durable
  ID lookup. Do not treat a GSearch path lookup as DAWA parity; live checks on
  2026-06-27 returned 404 for `/husnummer/{id}`.
- `adgangsadresser_reverse`: **No direct GSearch equivalent**. For
  nearest-address behavior, a **custom GSearch spatial lookup** can combine
  EPSG:25832 `INTERSECTS`/`DWITHIN` filters with nearest-result selection. Use
  a documented reverse/spatial service or another authoritative spatial lookup
  when exact DAWA reverse semantics matter.
- `historik_adgangsadresser`: **No GSearch equivalent**. Use authoritative
  history data where available.
- `datavask_adgangsadresser`: **Adressevask** when available, or a separate
  normalization flow. The official Adressevask replacement is not available yet
  as of 2026-06-27 and is expected at the end of August 2026.

### `adresser`

Generated operations:

```txt
adresser
adresser_autocomplete
adresser_{id}
historik_adresser
datavask_adresser
```

Migration path:

- `adresser_autocomplete`: **GSearch direct** via `adresse`.
- `adresser` text search: **GSearch partial** via `adresse` for best-match
  search; use DAR/Datafordeler for full query semantics.
- `adresser_{id}`: **DAR/Datafordeler/authoritative service** for durable ID
  lookup. Do not treat a GSearch path lookup as DAWA parity; live checks on
  2026-06-27 returned 404 for `/adresse/{id}`.
- `historik_adresser`: **No GSearch equivalent**. Use authoritative history data
  where available.
- `datavask_adresser`: **Adressevask** when available, or a separate
  normalization flow. The official Adressevask replacement is not available yet
  as of 2026-06-27 and is expected at the end of August 2026.

### Generic DAWA autocomplete

Generated operation:

```txt
autocomplete
```

Migration path:

- Use **Adressevælger** when you want SDFI's official picker/successor to the
  DAWA autocomplete value proposition.
- Use **GSearch direct** when you are building your own autocomplete UI. Query
  `husnummer` and `adresse` for address selection and other resources as needed.

## Roads, road names, and postcodes

### Named roads, road names, and road sections

Generated operations:

```txt
navngivneveje
navngivneveje_autocomplete
navngivneveje_{id}
navngivneveje_{id}_naboer
vejnavne
vejnavne_autocomplete
vejnavne_{navn}
vejstykker
vejstykker_autocomplete
vejstykker_reverse
vejstykker_{kommunekode}_{kode}
vejstykker_{kommunekode}_{kode}_naboer
```

Migration path:

- `*_autocomplete` and road text search: **GSearch direct/partial** via
  `navngivenvej`.
- ID/name/code lookup and neighbor endpoints:
  **DAR/Datafordeler/authoritative service**.
- `vejstykker_reverse`: **No direct GSearch equivalent**. Use
  DAR/Datafordeler/authoritative spatial data. A custom spatial index or
  GSearch spatial filter can support a product-specific nearest-road workflow,
  but that is not DAWA endpoint parity.

### Road-name postcode relations

Generated operations:

```txt
vejnavnpostnummerrelationer
vejnavnpostnummerrelationer_autocomplete
vejnavnpostnummerrelationer_{postnr}_{vejnavn}
```

Migration path:

- For user-facing road/postcode search, combine **GSearch** `navngivenvej` and
  `postnummer` where that is enough.
- For the relation object itself, use **DAR/Datafordeler/authoritative service**.
  GSearch does not expose a DAWA-equivalent relation endpoint.

### Postcodes

Generated operations:

```txt
postnumre
postnumre_autocomplete
postnumre_reverse
postnumre_{nr}
```

Migration path:

- `postnumre_autocomplete` and postcode search: **GSearch direct** via
  `postnummer`.
- `postnumre_{nr}` lookup: **DAR/Datafordeler/DAGI** or another authoritative
  dataset path.
- `postnumre_reverse`: **No direct GSearch equivalent**. You can use a spatial
  lookup over postcode polygons; for GSearch-assisted implementations, use
  EPSG:25832 ECQL filters and verify the behavior against your product's
  containment rules.

### Supplementary city names

Generated operations:

```txt
supplerendebynavne
supplerendebynavne_autocomplete
supplerendebynavne_{navn}
supplerendebynavne2
supplerendebynavne2_autocomplete
supplerendebynavne2_reverse
supplerendebynavne2_{dagi_id}
```

Migration path:

- **No direct GSearch resource**. GSearch address/house-number results can
  include `supplerendebynavn`, but standalone supplementary-city-name search,
  lookup, and reverse workflows need **DAR/Datafordeler/DAGI** or another
  authoritative source.

## DAGI administrative geography

### GSearch-supported DAGI resources

Generated operations:

```txt
kommuner
kommuner_autocomplete
kommuner_reverse
kommuner_{kode}
opstillingskredse
opstillingskredse_autocomplete
opstillingskredse_reverse
opstillingskredse_{kode}
politikredse
politikredse_autocomplete
politikredse_reverse
politikredse_{kode}
regioner
regioner_autocomplete
regioner_reverse
regioner_{kode}
retskredse
retskredse_autocomplete
retskredse_reverse
retskredse_{kode}
sogne
sogne_autocomplete
sogne_reverse
sogne_{kode}
```

Migration path:

- `*_autocomplete` and name/code-ish search: **GSearch direct** via singular
  resources `kommune`, `opstillingskreds`, `politikreds`, `region`,
  `retskreds`, and `sogn`.
- `*_{kode}` lookup and `*_reverse`: **DAR/Datafordeler/DAGI** or spatial
  lookup. GSearch can help users find a result, but it is not a general reverse
  or full-object lookup API.

### DAGI resources without a default GSearch resource

Generated operations:

```txt
afstemningsomraader
afstemningsomraader_autocomplete
afstemningsomraader_reverse
afstemningsomraader_{kommunekode}_{nummer}
landsdele
landsdele_autocomplete
landsdele_reverse
landsdele_{nuts3}
menighedsraadsafstemningsomraader
menighedsraadsafstemningsomraader_autocomplete
menighedsraadsafstemningsomraader_reverse
menighedsraadsafstemningsomraader_{kommunekode}_{nummer}
storkredse
storkredse_autocomplete
storkredse_reverse
storkredse_{nummer}
valglandsdele
valglandsdele_autocomplete
valglandsdele_reverse
valglandsdele_{bogstav}
```

Migration path:

- **No default GSearch resource** for these DAWA endpoint families.
- Use **DAR/Datafordeler/DAGI**, OGC/WFS/download services, or another authoritative
  dataset path.
- If you only need a user-facing free-text search across broader geodata, check
  whether `stednavn` or a supported DAGI resource meets the product need, but do
  not describe that as endpoint parity.

## Cadastral endpoints

### Parcels and cadastral districts

Generated operations:

```txt
ejerlav
ejerlav_autocomplete
ejerlav_{kode}
jordstykker
jordstykker_autocomplete
jordstykker_reverse
jordstykker_{ejerlavkode}_{matrikelnr}
```

Migration path:

- `jordstykker_autocomplete` and parcel search: **GSearch partial** via
  `matrikel` and, when relevant, `matrikel_udgaaet`.
- `ejerlav` search/autocomplete: **No direct GSearch resource**. GSearch
  `matrikel` can be filtered by fields such as `ejerlavskode`, but standalone
  cadastral-district lookup belongs in **Datafordeler/Matriklen**.
- parcel lookup and reverse: **Datafordeler/Matriklen**, OGC/WFS, or spatial
  lookup. GSearch is not a DAWA reverse endpoint, though `matrikel` spatial
  filters can help product-specific map search workflows.

## Place names and settlements

Generated operations:

```txt
bebyggelser
bebyggelser_{id}
steder
steder_{id}
stednavne
stednavne_autocomplete
stednavne_{id}
stednavne2
stednavne2_autocomplete
stednavne2_{sted_id}_{navn}
stednavntyper
stednavntyper_{hovedtype}
```

Migration path:

- `stednavne*_autocomplete` and place-name search: **GSearch direct/partial**
  via `stednavn`.
- `bebyggelser`, `steder`, `stednavntyper`, ID lookup, and type lookup:
  **DAR/Datafordeler/authoritative service** or the Danish place-name dataset.
  GSearch can support search UX, not full DAWA object parity.

## BBR, OIS, and BBR light endpoints

Generated operations:

```txt
bbrlight_bygninger
bbrlight_bygninger_{id}
bbrlight_bygningspunkter
bbrlight_bygningspunkter_{id}
bbrlight_ejerskaber
bbrlight_ejerskaber_{id}
bbrlight_enheder
bbrlight_enheder_{id}
bbrlight_etager
bbrlight_etager_{id}
bbrlight_grunde
bbrlight_grunde_{id}
bbrlight_kommuner
bbrlight_kommuner_{id}
bbrlight_matrikelreferencer
bbrlight_opgange
bbrlight_opgange_{id}
bbrlight_tekniskeanlaeg
bbrlight_tekniskeanlaeg_{id}
ois_bygninger
ois_bygninger_{id}
ois_bygningspunkter
ois_bygningspunkter_{id}
ois_ejerskaber
ois_ejerskaber_{id}
ois_enheder
ois_enheder_{id}
ois_etager
ois_etager_{id}
ois_grunde
ois_grunde_{id}
ois_kommuner
ois_kommuner_{id}
ois_matrikelreferencer
ois_opgange
ois_opgange_{id}
ois_tekniskeanlaeg
ois_tekniskeanlaeg_{id}
```

Migration path:

- **No GSearch equivalent for BBR/OIS data.**
- DAWA's own documentation states BBR data on DAWA has been phased out.
- For address-led workflows, use GSearch to select and persist a `husnummer`
  UUID, then resolve it through DAR BFE/Datafordeler to `jordstykkeLokalId` and
  BFE number before querying BBR/MAT or the current authoritative
  property/register service required by your use case.
- For direct BBR/OIS object access, use Datafordeler BBR/MAT or the current
  authoritative data access path required by your use case.

## Building polygons

Generated operations:

```txt
bygninger
bygninger_{id}
```

Migration path:

- **No GSearch equivalent for building footprints.**
- For targeted lookup near an address, use **Datafordeler GeoDanmark Vektor
  GraphQL** `GEODKV_Bygning` with a spatial filter on `geometri`.
- For bulk/local-copy/GIS workflows, use **GeoDanmark Vektor Fildownload** or
  **GeoDanmark Vektor WFS entities**.
- If the workflow starts from an address, use GSearch only to get coordinates
  and DAR IDs, then switch to the building-footprint source. Avoid new
  dependencies on legacy WFS/FTP filudtræk paths without checking the current
  Datafordeler data overview and phase-out notices.

## Replication and local-copy endpoints

### Current replication API

Generated operations:

```txt
replikering_udtraek
replikering_haendelser
replikering_transaktioner
replikering_senestetransaktion
```

Migration path:

- **No GSearch equivalent.**
- Use **Datafordeler events**, Dataforsyningen downloads, OGC/WFS services, or
  dataset-specific replication/export mechanisms. GSearch is an online search
  API, not a change-feed or local-copy API.

### Deprecated replication API

Generated operations:

```txt
replikering_senesteSekvensnummer
replikering_adgangsadresser
replikering_adgangsadresser_haendelser
replikering_adresser
replikering_adresser_haendelser
replikering_afstemningsområdetilknytninger
replikering_afstemningsområdetilknytninger_haendelser
replikering_ejerlav
replikering_ejerlav_haendelser
replikering_jordstykketilknytninger
replikering_jordstykketilknytninger_haendelser
replikering_kommunetilknytninger
replikering_kommunetilknytninger_haendelser
replikering_landsdelstilknytninger
replikering_landsdelstilknytninger_haendelser
replikering_menighedsrådsafstemningsområdetilknytninger
replikering_menighedsrådsafstemningsområdetilknytninger_haendelser
replikering_opstillingskredstilknytninger
replikering_opstillingskredstilknytninger_haendelser
replikering_politikredstilknytninger
replikering_politikredstilknytninger_haendelser
replikering_postnummertilknytninger
replikering_postnummertilknytninger_haendelser
replikering_postnumre
replikering_postnumre_haendelser
replikering_regionstilknytninger
replikering_regionstilknytninger_haendelser
replikering_retskredstilknytninger
replikering_retskredstilknytninger_haendelser
replikering_sognetilknytninger
replikering_sognetilknytninger_haendelser
replikering_stednavntilknytninger
replikering_stednavntilknytninger_haendelser
replikering_storkredstilknytninger
replikering_storkredstilknytninger_haendelser
replikering_supplerendebynavntilknytninger
replikering_supplerendebynavntilknytninger_haendelser
replikering_valglandsdelstilknytninger
replikering_valglandsdelstilknytninger_haendelser
replikering_vejstykkepostnummerrelationer
replikering_vejstykkepostnummerrelationer_haendelser
replikering_vejstykker
replikering_vejstykker_haendelser
replikering_zonetilknytninger
replikering_zonetilknytninger_haendelser
```

Migration path:

- **No GSearch equivalent.**
- Treat these as local-copy/change-feed workflows. Use Datafordeler events,
  Dataforsyningen downloads, OGC/WFS services, or dataset-specific exports.

## Experimental/historical endpoint

Generated operation:

```txt
darhistorik
```

Migration path:

- **No GSearch equivalent.**
- Use an authoritative DAR history source if available for the specific use
  case. Do not migrate history queries to GSearch.

## Coverage assertion

Every generated DAWA operation listed above has one of these outcomes:

- direct/partial GSearch path for search/autocomplete;
- Adressevælger or the planned Adressevask service for address-picker/data-washing workflows;
- custom GSearch spatial lookup where nearest-object behavior is enough and
  endpoint parity is not required;
- DAR/Datafordeler/authoritative service path for lookups, register data,
  reverse, history, replication, BBR/OIS, BFE, and bulk/local-copy workflows;
- explicit "no GSearch equivalent" where endpoint parity does not exist.

The absence of a GSearch path is intentional for non-search features. The
migration goal is to stop using DAWA, not to force all DAWA behavior through
GSearch.
