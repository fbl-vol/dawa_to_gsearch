# Implementation patterns

These patterns are for applications that currently lean on DAWA for more than
simple autocomplete. They keep GSearch in its lane as a search API and use
DAR/Datafordeler when the workflow needs authoritative register data.

## Address autocomplete with downstream IDs

For an address picker, query both GSearch address resources:

- `husnummer` for building-level/access-address results;
- `adresse` for unit-level results with floor and door.

Merge the two result sets in your own application model. Prefer the `adresse`
result for display when the user selected a unit, but preserve the parent
`husnummer` UUID for workflows that later need BBR, BFE, or parcel data.

A practical matching key is:

```txt
kommunekode + vejkode + husnummer
```

This gives you an internal selected-address shape like:

```ts
type SelectedAddress = {
  provider: "dataforsyningen-gsearch";
  resource: "adresse" | "husnummer";
  id: string;          // adresse UUID for unit selections, husnummer UUID otherwise
  husnummerId?: string; // parent husnummer UUID for BBR/DAR/BFE lookups
  label: string;
  kommunekode?: string;
  vejkode?: string;
  postnummer?: string;
  longitude?: number;
  latitude?: number;
};
```

Do not store only `visningstekst`. Store the UUIDs and structured fields.

## Address to BFE, parcel, and BBR

GSearch does not replace BBR or full DAR object lookup. Use it to get the right
DAR UUID, then switch to Datafordeler/DAR/BBR.

One concrete migration flow:

1. Search GSearch `husnummer`/`adresse`.
2. Keep the `husnummer` UUID.
3. Resolve the `husnummer` UUID through DAR BFE to get `jordstykkeLokalId` and
   `samletFastEjendom` (BFE number).
4. Use `jordstykkeLokalId`, BFE number, or BBR building IDs with
   Datafordeler BBR/MAT/DAR GraphQL or the current Datafordeler service required
   by your use case.

Legacy REST examples seen in existing integrations use:

```txt
https://services.datafordeler.dk/DAR/DAR_BFE_Public/1/REST/husnummerTilBygningBfe
```

with parameters such as:

```txt
HusnummerId=<husnummer-uuid>
format=json
username=<datafordeler-username>
password=<DATAFORDELER_PASSWORD>
```

Datafordeler is transitioning REST services toward GraphQL/API-key/OAuth based
access. For new work, check the current Datafordeler DAR, BBR, MAT, and
transition documentation before choosing REST or GraphQL.

Try the bridge shape with:

```sh
examples/curl/dar-bfe-from-husnummer.sh "<husnummer-uuid>"
```

## Reverse geocoding without a DAWA reverse endpoint

GSearch has no DAWA-style `/reverse` endpoint. If the product only needs
"nearest address for coordinates", you can build a custom spatial lookup with
GSearch filters:

1. Convert WGS84 latitude/longitude to EPSG:25832 (ETRS89/UTM32N).
2. Find the containing postal district with `postnummer` and an ECQL spatial
   filter:

```txt
INTERSECTS(geometri,POINT(<easting> <northing>))
```

3. Search nearby `husnummer` results with increasing radii:

```txt
DWITHIN(geometri,POINT(<easting> <northing>),100,meters)
DWITHIN(geometri,POINT(<easting> <northing>),300,meters)
DWITHIN(geometri,POINT(<easting> <northing>),1000,meters)
```

4. Compare returned native EPSG:25832 coordinates and choose the nearest
   `husnummer`.
5. If the result feeds register workflows, resolve the chosen `husnummer` UUID
   through DAR/Datafordeler before using it as authoritative BBR/BFE input.

Try the default Gedser example with:

```sh
examples/curl/spatial-nearest-husnummer.sh
```

This is a migration pattern, not endpoint parity. If the old DAWA call required
official historical behavior, specific containment rules, or administrative
reverse lookups, use Datafordeler, OGC/WFS, downloads, or a purpose-built
spatial index instead.

## Building and property workflows

For building or property pages, use GSearch only for the user-facing address
search step. After selection:

- resolve `husnummer` to BFE/jordstykke through DAR/Datafordeler;
- fetch building, unit, floor, entrance, ground, and technical-installation
  data from BBR/Datafordeler;
- use MAT/Datafordeler for property and cadastral attributes;
- use Datafordeler GeoDanmark Vektor GraphQL `GEODKV_Bygning` for targeted
  building-footprint lookup near a selected address;
- use GeoDanmark Vektor Fildownload or WFS entities for bulk/local-copy/GIS
  workflows.

This split prevents a common migration mistake: treating a search result as a
complete register object.

For an on-demand building footprint lookup, query Datafordeler GraphQL register
`GEODKV` version `v2`, entity `GEODKV_Bygning`, with a spatial filter on
`geometri`:

```graphql
query GetBygninger(
  $where: GEODKV_BygningFilterInput
  $registreringstid: DafDateTime
  $virkningstid: DafDateTime
) {
  GEODKV_Bygning(
    first: 100
    where: $where
    registreringstid: $registreringstid
    virkningstid: $virkningstid
  ) {
    nodes {
      geometri { crs wkt }
      bygningstype
      status
      id_lokalId
      BBRUUID
    }
  }
}
```

Use an EPSG:25832 WKT geometry in the variables. Use current UTC instants for
`registreringstid` and `virkningstid`; the shell example sets them
automatically.

```json
{
  "first": 100,
  "registreringstid": "2026-06-27T00:00:00Z",
  "virkningstid": "2026-06-27T00:00:00Z",
  "where": {
    "geometri": {
      "intersects": {
        "crs": 25832,
        "wkt": "POLYGON((689245 6051777, 689265 6051777, 689265 6051797, 689245 6051797, 689245 6051777))"
      }
    }
  }
}
```

Try the credential-safe example:

```sh
examples/curl/geodkv-building-footprints.sh
```

The older GeoDanmark Vektor WFS page reviewed on 2026-06-27 marks that specific
WFS service for phase-out in 2026. The old GeoDanmark `filudtræk` documentation
also points users toward Fildownload. Do not hard-code a legacy WFS or FTP
filudtræk path for new long-lived work without checking the current
Datafordeler data overview.

## Data washing, history, and replication

These DAWA features should not be reimplemented with GSearch:

- `datavask`: use Klimadatastyrelsen/SDFI's planned Adressevask replacement
  when available, or a dedicated normalization/validation pipeline. As of
  2026-06-27, the official docs say the new Adressevask API is not available
  yet and is expected at the end of August 2026.
- history: use DAR/Datafordeler historical data where available for the needed
  entity and time semantics.
- replication/change feeds: use Datafordeler events, file downloads, WFS, or a
  dataset-specific export path.

GSearch can help a user find an object. It should not be your local-copy,
history, or normalization engine.

For copy-paste commands across all resource types, see
[examples gallery](examples-gallery.md).

For Datafordeler-specific BFE, BBR/MAT, GraphQL, events, downloads, and WFS/OGC
guidance, see [DAR and Datafordeler guide](datafordeler-dar-guide.md).
