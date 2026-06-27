#!/usr/bin/env python3
"""Query Datafordeler GeoDanmark Vektor building footprints."""

from __future__ import annotations

import argparse
import datetime as dt
import json
import os
import sys
import urllib.parse
import urllib.request
from typing import Any, Optional


ENDPOINT = "https://graphql.datafordeler.dk/GEODKV/v2"

QUERY = """query GetBygninger($first: Int, $where: GEODKV_BygningFilterInput, $registreringstid: DafDateTime, $virkningstid: DafDateTime) {
  GEODKV_Bygning(first: $first, where: $where, registreringstid: $registreringstid, virkningstid: $virkningstid) {
    nodes {
      geometri { crs wkt }
      bygningstype
      status
      id_lokalId
      BBRUUID
    }
  }
}"""


def geodkv_building_footprints(
    *,
    easting: int = 689255,
    northing: int = 6051787,
    bbox_size_meters: int = 20,
    limit: int = 10,
    api_key: Optional[str] = None,
    timestamp: Optional[str] = None,
    timeout: float = 10,
) -> dict[str, Any]:
    api_key = api_key or os.environ.get("DATAFORDELEREN_API_KEY") or os.environ.get("DATAFORDELER_GRAPHQL_API_KEY")
    if not api_key:
        raise RuntimeError("DATAFORDELEREN_API_KEY or DATAFORDELER_GRAPHQL_API_KEY is required")

    timestamp = timestamp or dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    half = bbox_size_meters // 2
    min_e = easting - half
    min_n = northing - half
    max_e = easting + half
    max_n = northing + half
    wkt = f"POLYGON(({min_e} {min_n}, {max_e} {min_n}, {max_e} {max_n}, {min_e} {max_n}, {min_e} {min_n}))"

    body = {
        "query": QUERY,
        "variables": {
            "first": limit,
            "registreringstid": timestamp,
            "virkningstid": timestamp,
            "where": {
                "geometri": {
                    "intersects": {
                        "crs": 25832,
                        "wkt": wkt,
                    }
                }
            },
        },
    }
    encoded_key = urllib.parse.quote(api_key, safe="")
    request = urllib.request.Request(
        f"{ENDPOINT}?apiKey={encoded_key}",
        data=json.dumps(body).encode("utf-8"),
        headers={
            "accept": "application/json",
            "content-type": "application/json",
        },
        method="POST",
    )

    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            payload = json.load(response)
    except urllib.error.HTTPError as error:
        body_text = error.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"GEODKV query failed: HTTP {error.code} {body_text}") from error
    except urllib.error.URLError as error:
        raise RuntimeError(f"GEODKV request failed: {error.reason}") from error

    if not isinstance(payload, dict):
        raise RuntimeError(f"GEODKV returned {type(payload).__name__}, expected object")
    return payload


def main() -> int:
    parser = argparse.ArgumentParser(description="Query GEODKV_Bygning around an EPSG:25832 point.")
    parser.add_argument("--easting", type=int, default=689255)
    parser.add_argument("--northing", type=int, default=6051787)
    parser.add_argument("--bbox-size-meters", type=int, default=20)
    parser.add_argument("--limit", type=int, default=10)
    args = parser.parse_args()

    payload = geodkv_building_footprints(
        easting=args.easting,
        northing=args.northing,
        bbox_size_meters=args.bbox_size_meters,
        limit=args.limit,
    )
    json.dump(payload, sys.stdout, ensure_ascii=False, indent=2)
    sys.stdout.write("\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
