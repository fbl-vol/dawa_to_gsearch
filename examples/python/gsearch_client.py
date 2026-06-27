#!/usr/bin/env python3
"""Small dependency-free GSearch client for address suggestions."""

from __future__ import annotations

import argparse
import json
import os
import sys
import urllib.parse
import urllib.request
from concurrent.futures import ThreadPoolExecutor
from typing import Any, Optional


BASE_URL = "https://api.dataforsyningen.dk/rest/gsearch/v2.0"


def search_gsearch(
    resource: str,
    query: str,
    *,
    token: Optional[str] = None,
    limit: int = 10,
    srid: int = 4326,
    timeout: float = 10,
) -> list[dict[str, Any]]:
    token = token or os.environ.get("GSEARCH_TOKEN")
    if not token:
        raise RuntimeError("GSEARCH_TOKEN is required")

    params = urllib.parse.urlencode(
        {
            "token": token,
            "q": query,
            "limit": str(limit),
            "srid": str(srid),
        }
    )
    url = f"{BASE_URL}/{resource}?{params}"
    request = urllib.request.Request(url, headers={"accept": "application/json"})

    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            payload = json.load(response)
    except urllib.error.HTTPError as error:
        body = error.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"GSearch {resource} failed: HTTP {error.code} {body}") from error
    except urllib.error.URLError as error:
        raise RuntimeError(f"GSearch {resource} request failed: {error.reason}") from error

    if not isinstance(payload, list):
        raise RuntimeError(f"GSearch {resource} returned {type(payload).__name__}, expected list")
    return payload


def search_address_suggestions(
    query: str,
    *,
    token: Optional[str] = None,
    limit: int = 10,
    srid: int = 4326,
    timeout: float = 10,
) -> list[dict[str, Any]]:
    with ThreadPoolExecutor(max_workers=2) as executor:
        husnummer_future = executor.submit(
            search_gsearch,
            "husnummer",
            query,
            token=token,
            limit=limit,
            srid=srid,
            timeout=timeout,
        )
        adresse_future = executor.submit(
            search_gsearch,
            "adresse",
            query,
            token=token,
            limit=limit,
            srid=srid,
            timeout=timeout,
        )
        return merge_address_results(husnummer_future.result(), adresse_future.result())


def merge_address_results(
    husnummer_results: list[dict[str, Any]],
    adresse_results: list[dict[str, Any]],
) -> list[dict[str, Any]]:
    husnummer_ids: dict[str, str] = {}
    for result in husnummer_results:
        key = address_key(result, "husnummertekst")
        result_id = result.get("id")
        if key and isinstance(result_id, str):
            husnummer_ids[key] = result_id

    seen: set[str] = set()
    merged: list[dict[str, Any]] = []

    for result in adresse_results:
        key = address_key(result, "husnummer")
        normalized = normalize_result("adresse", result, husnummer_ids.get(key or ""))
        label = normalized.get("label")
        if isinstance(label, str) and label not in seen:
            seen.add(label)
            merged.append(normalized)

    for result in husnummer_results:
        result_id = result.get("id")
        normalized = normalize_result("husnummer", result, result_id if isinstance(result_id, str) else None)
        label = normalized.get("label")
        if isinstance(label, str) and label not in seen:
            seen.add(label)
            merged.append(normalized)

    return merged


def address_key(result: dict[str, Any], house_number_field: str) -> Optional[str]:
    parts = [
        result.get("kommunekode"),
        result.get("vejkode"),
        result.get(house_number_field),
    ]
    if any(part in (None, "") for part in parts):
        return None
    return ":".join(str(part) for part in parts)


def normalize_result(
    resource: str,
    result: dict[str, Any],
    husnummer_id: Optional[str],
) -> dict[str, Any]:
    lon_lat = first_lon_lat(result.get("geometri"))
    normalized: dict[str, Any] = {
        "provider": "dataforsyningen-gsearch",
        "resource": resource,
        "id": result.get("id"),
        "husnummerId": husnummer_id,
        "label": result.get("visningstekst"),
        "kommunekode": result.get("kommunekode"),
        "vejkode": result.get("vejkode"),
        "postnummer": result.get("postnummer"),
        "raw": result,
    }
    if lon_lat:
        normalized["longitude"] = lon_lat[0]
        normalized["latitude"] = lon_lat[1]
    return normalized


def first_lon_lat(geometry: Any) -> Optional[tuple[float, float]]:
    if not isinstance(geometry, dict):
        return None
    return first_position(geometry.get("coordinates"))


def first_position(value: Any) -> Optional[tuple[float, float]]:
    if not isinstance(value, list):
        return None
    if len(value) >= 2 and isinstance(value[0], (int, float)) and isinstance(value[1], (int, float)):
        return float(value[0]), float(value[1])
    for child in value:
        position = first_position(child)
        if position:
            return position
    return None


def main() -> int:
    parser = argparse.ArgumentParser(description="Search GSearch adresse and husnummer in parallel.")
    parser.add_argument("query", nargs="?", default="Søbakkevej 8, Tilst")
    parser.add_argument("--limit", type=int, default=5)
    parser.add_argument("--srid", type=int, default=4326)
    args = parser.parse_args()

    suggestions = search_address_suggestions(args.query, limit=args.limit, srid=args.srid)
    json.dump(suggestions, sys.stdout, ensure_ascii=False, indent=2)
    sys.stdout.write("\n")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
