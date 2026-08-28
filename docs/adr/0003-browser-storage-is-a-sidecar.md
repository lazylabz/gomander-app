# Browser storage is a private sidecar, not a port

`store/sidebarSections.ts` reads and writes `localStorage` directly, while every
other external dependency in the frontend sits behind a port with a fake adapter
(see ADR-0001).

The inconsistency is deliberate. A port earns its place when it has more than one
real implementation, or when tests cannot run against the real thing. Neither
holds: the jsdom test environment provides browser storage, so a port here would
have exactly one real adapter and one fake that behaves identically. The storage
key, the JSON encoding and the `try`/`catch` around a rejected read stay
implementation details of this one module, and callers see only
`isSidebarSectionOpen` / `setSidebarSectionOpen`.

The line to hold: this applies to per-viewer convenience state. State the app
would be wrong without does not belong in browser storage at all.
