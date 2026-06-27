# GSearch response shapes

GSearch responses are JSON arrays. The exact fields depend on the resource, but
all search results should be handled as structured objects, not strings.

## Common fields

| Field | Meaning |
| --- | --- |
| `id` | Authoritative UUID for the result object. |
| `visningstekst` | Human-readable display text for suggestion lists. |
| `geometri` | GeoJSON geometry. |
| `bbox` | Bounding box on resources where it is available. |

The official GSearch docs state that object geometry is included as GeoJSON and
that the default/native spatial reference is EPSG:25832. GSearch UI also exposes
`srid` values including `4326`. For web apps, request `srid=4326` if you need
longitude/latitude output.

Spatial filters still use EPSG:25832.

## `husnummer`

Use `husnummer` for building-level/access-address results.

Important fields commonly used by applications:

```ts
type GSearchHusnummer = {
  id: string;
  kommunekode?: string;
  kommunenavn?: string;
  vejkode?: string;
  vejnavn?: string;
  husnummertekst?: string;
  postnummer?: string;
  postnummernavn?: string;
  supplerendebynavn?: string | null;
  visningstekst: string;
  geometri?: GeoJSON.Geometry;
  vejpunkt_geometri?: GeoJSON.Geometry;
};
```

Persist `id` as the selected building/access-address UUID.

## `adresse`

Use `adresse` for unit-level results such as floor and door.

```ts
type GSearchAdresse = {
  id: string;
  kommunekode?: string;
  kommunenavn?: string;
  vejkode?: string;
  vejnavn?: string;
  husnummer?: string;
  etagebetegnelse?: string | null;
  doerbetegnelse?: string | null;
  postnummer?: string;
  postnummernavn?: string;
  supplerendebynavn?: string | null;
  visningstekst: string;
  geometri?: GeoJSON.Geometry;
  vejpunkt_geometri?: GeoJSON.Geometry;
};
```

If your downstream workflow needs a building/access-address UUID, query
`husnummer` in parallel and match by `kommunekode + vejkode + husnummer`. Use
that parent `husnummer` UUID for DAR BFE, BBR, parcel, and property workflows.

## Unified selected result

A good internal model is narrower than the raw GSearch response:

```ts
type DanishAddressSelection = {
  provider: "dataforsyningen-gsearch";
  resource: "adresse" | "husnummer";
  id: string;
  husnummerId?: string;
  label: string;
  kommunekode?: string;
  vejkode?: string;
  postnummer?: string;
  longitude?: number;
  latitude?: number;
  raw?: unknown;
};
```

Keep `raw` optional. It is useful during migration, but long-term app code should
depend on the fields it actually needs.

## IDs for register joins

Do not assume one UUID is enough for every downstream system:

- `adresse.id` identifies the unit-level DAR address.
- `husnummer.id` identifies the access address/house number.
- BFE/BBR workflows commonly need the `husnummer` UUID first, then a
  DAR/Datafordeler lookup to resolve BFE number, `jordstykke`, building, unit,
  floor, entrance, or ground data.

If your UI displays `adresse` suggestions first, add your own `husnummerId`
field by matching the parallel `husnummer` results structurally.

## Coordinates

When you request `srid=4326`, GeoJSON coordinates are longitude/latitude:

```ts
function firstLonLat(geometry: { coordinates?: unknown }): [number, number] | null {
  return firstPosition(geometry.coordinates);
}

function firstPosition(value: unknown): [number, number] | null {
  if (!Array.isArray(value)) {
    return null;
  }
  if (typeof value[0] === "number" && typeof value[1] === "number") {
    return [value[0], value[1]];
  }
  for (const child of value) {
    const position = firstPosition(child);
    if (position) {
      return position;
    }
  }
  return null;
}
```

Do not assume every resource uses the same geometry type. GeoJSON `Point`
coordinates are a single `[longitude, latitude]` pair, `MultiPoint` and
`LineString` use arrays of pairs, and polygons nest the pairs more deeply.
Address-like resources often behave like point/multipoint data, while DAGI and
place-name resources can have polygon or multipolygon geometry.

## Display text

Use `visningstekst` for UI display. Do not parse it to recover road names,
postcodes, floor, or door. Use the structured fields returned by the endpoint.

Display text can change without your app's data model changing. UUIDs and
structured fields are the stable integration surface.
