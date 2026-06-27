const BASE_URL = "https://api.dataforsyningen.dk/rest/gsearch/v2.0";

export async function searchGSearch(resource, query, options = {}) {
  const {
    token = defaultToken(),
    limit = 10,
    srid = 4326,
    signal,
  } = options;

  if (!token) {
    throw new Error("GSEARCH_TOKEN is required");
  }

  const url = new URL(`${BASE_URL}/${resource}`);
  url.searchParams.set("token", token);
  url.searchParams.set("q", query);
  url.searchParams.set("limit", String(limit));
  url.searchParams.set("srid", String(srid));

  const response = await fetch(url, {
    headers: { Accept: "application/json" },
    signal,
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`GSearch ${resource} failed: HTTP ${response.status} ${body}`);
  }

  return response.json();
}

function defaultToken() {
  if (typeof process !== "undefined" && process.env?.GSEARCH_TOKEN) {
    return process.env.GSEARCH_TOKEN;
  }
  return undefined;
}

export async function searchAddressSuggestions(query, options = {}) {
  const [husnummer, adresse] = await Promise.all([
    searchGSearch("husnummer", query, options),
    searchGSearch("adresse", query, options),
  ]);

  return mergeAddressResults(husnummer, adresse);
}

function mergeAddressResults(husnummerResults, adresseResults) {
  const husnummerIds = new Map();
  for (const result of husnummerResults) {
    const key = [
      result.kommunekode,
      result.vejkode,
      result.husnummertekst,
    ].join(":");
    husnummerIds.set(key, result.id);
  }

  const seen = new Set();
  const merged = [];

  for (const result of adresseResults) {
    const key = [result.kommunekode, result.vejkode, result.husnummer].join(":");
    const normalized = normalizeResult("adresse", result, husnummerIds.get(key));
    if (!seen.has(normalized.label)) {
      seen.add(normalized.label);
      merged.push(normalized);
    }
  }

  for (const result of husnummerResults) {
    const normalized = normalizeResult("husnummer", result, result.id);
    if (!seen.has(normalized.label)) {
      seen.add(normalized.label);
      merged.push(normalized);
    }
  }

  return merged;
}

function normalizeResult(resource, result, husnummerId) {
  const [longitude, latitude] = firstLonLat(result.geometri) ?? [];
  return {
    provider: "dataforsyningen-gsearch",
    resource,
    id: result.id,
    husnummerId,
    label: result.visningstekst,
    kommunekode: result.kommunekode,
    vejkode: result.vejkode,
    postnummer: result.postnummer,
    longitude,
    latitude,
    raw: result,
  };
}

function firstLonLat(geometry) {
  return firstPosition(geometry?.coordinates);
}

function firstPosition(value) {
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
