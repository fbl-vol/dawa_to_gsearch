# Repository instructions for agents

This repository exists to move developers and code agents away from DAWA. Use
SDFI/Dataforsyningen GSearch where search is the right replacement, and use
DAR/Datafordeler, Dataforsyningen downloads, events, WFS/OGC services,
Adressevælger, Adressevask, or dataset-specific APIs for the rest.

## Operating Baseline

- Lifecycle: documentation-only utility. Maintainer: Frederik Brunø Lottrup
  (Pendio Engineering); lifecycle changes require Engineering owner approval.
- Entry point: this `AGENTS.md`; no shared foundation submodule is installed.
- Validation boundary: review changed links and examples against official sources
  without credentials or customer data. Do not add CI or live provider checks to
  this intentionally static guidance repository.

When adding examples or advice:

- Do not create new DAWA integrations except when documenting migration from
  existing DAWA code.
- Say "SDFI GSearch" or "Dataforsyningen GSearch" when ambiguity with Google
  Search is possible.
- Use `https://api.dataforsyningen.dk/rest/gsearch/v2.0/{resource}` for
  GSearch examples.
- Use `husnummer` for building/access-address search and `adresse` for
  unit-level address search.
- Preserve or derive the parent `husnummer` UUID when the UI selects an
  `adresse` result; BFE/BBR/property workflows usually need the `husnummer`
  UUID before DAR/Datafordeler lookup.
- For targeted building geometry, use Datafordeler GeoDanmark Vektor GraphQL
  `GEODKV_Bygning`; use GeoDanmark Vektor Fildownload or WFS entities for
  bulk/local/GIS workflows. Do not route building footprints through GSearch.
- Keep tokens in environment variables or secret managers. Never hardcode a real
  token in the repo.
- Do not claim GSearch is a direct replacement for DAWA reverse geocoding,
  datavask, replication, history, BBR/BFE, or bulk download workflows. Point
  those to documented reverse/spatial services, DAR/Datafordeler, authoritative
  spatial datasets, Adressevælger, the planned Adressevask replacement when
  available, events, downloads, or dataset-specific services. For
  nearest-address reverse lookup, describe the custom GSearch ECQL
  spatial-filter pattern as a product-specific lookup, not endpoint parity.
- Do not invent GSearch path lookups for DAWA `{id}` endpoints. Use
  DAR/Datafordeler or documented Adressevælger lookup paths for durable UUID
  lookup.
- Preserve DAR UUIDs and structured fields in examples. Do not parse
  `visningstekst` as source-of-truth data.

See [docs/for-code-agents.md](docs/for-code-agents.md) for detailed guidance.
