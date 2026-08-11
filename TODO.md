# Implementation TODO

## Milestones (from ARCHITECTURE.md)

| # | Milestone | Status |
|---|---|---|
| 1 | MVP — `.container` files only, priority-1 field mappings | ✅ Done |
| 2 | Full compose parity — `.network`, `.volume`, `.image`, `.build`, all P1+P2 | ✅ Done |
| 3 | Opinionated defaults — all comquad transforms ported | ✅ Done |
| 4 | Deploy + systemd — `deploy.resources`, `deploy.restart_policy` | ✅ Done |
| 5 | Secrets + builds — `secrets:` and `build:` handled natively | ✅ Done |
| 6 | Integration — comquad imports the library, drops podlet dependency | ❌ |
| 7 | Deprecate podlet — comquad no longer requires podlet binary at runtime | ❌ |

## Testing Gaps (from ARCHITECTURE.md)

### Tier 0 — Serialization

| # | Item | Description |
|---|---|---|
| T0-1 | Golden files | ✅ Done — `testdata/serialization/` with `.golden` files for container, network, volume, image, build |
| T0-2 | Round-trip completeness | ✅ Done — round-trip tests for `.network`, `.volume`, `.image`, `.build` + empty directive round-trip |

### Tier 2 — Version Matrix

| # | Item | Description |
|---|---|---|
| T2-1 | Version boundary tests | ✅ Done — `TestVersion_Entrypoint_P3toP1`, `TestVersion_StopSignal_Gate`, `TestVersion_ExtraHosts_Gate` |
| T2-2 | Feature availability gates | ✅ Done — `TestVersion_NetworkAliases_Gate`, `TestVersion_LogOptions_Gate`, `TestVersion_Build_Available` |
| T2-3 | Section-switching fields | ✅ Done — `TestVersion_MemorySectionSwitch` |
| T2-4 | Fatal structural blocks | ✅ Done — `TestVersion_Build_FatalError` |
| T2-5 | Warning collection | ✅ Done — `TestVersion_WarningCollectionSmoke` + mapper-level warning tests |

### Tier 3 — Pipeline Integration

| # | Item | Description |
|---|---|---|
| T3-1 | Integration tests | ✅ Done — `transpile_test.go` with fixture YAML files in `testdata/` |
| T3-2 | Service combos | ✅ Done — `TestTranspile_SimpleWeb`, `TestTranspile_MultiService`, `TestTranspile_TopLevelOnly` |
| T3-3 | Option combinatorics | ✅ Done — `TestTranspile_OptionCombinatorics` (install+autoupdate, selinux+default network, port offset, default prefix) |
| T3-4 | Warning verification | ✅ Done — mapper-level `TestContainer_VersionGatedWarnings` + pipeline `TestVersion_*` tests |
| T3-5 | Edge cases | ✅ Done — `TestTranspile_EdgeCases_BuildService`, `TestTranspile_ExternalVolumesSkipped` |

### Tier 4 — End-to-End

Deferred to comquad's `tests/integration/` harness. The library itself does not start podman or systemd.

## Code Quality Issues (from code review)

### Critical

| # | Issue | Status |
|---|---|---|
| 279 | Non-deterministic map iteration | ✅ Fixed — sorted keys via `sortedKeys()` helper |
| 280 | Build + Image logic bug | ✅ Fixed — added Image fallback for < 5.2 |
| 281 | Conflicting CPU directives | ✅ Fixed — `svc.CPUs` takes priority over `svc.CPUQuota` and `deploy` |
| 282 | Conflicting restart directives | ✅ Fixed — `Deploy.RestartPolicy` takes priority over `svc.Restart` |

### Medium

| # | Issue | Status |
|---|---|---|
| 283 | Incomplete port offset logic | ✅ Fixed — rewritten to strip protocol first, then offset |
| 284 | SELinux false positives | ✅ Fixed — `hasSELinuxContext()` checks last colon segment only |
| 285 | Missing opinionated tests | ✅ Fixed — 25 tests covering all transforms |
| 286 | Hardcoded retry values | ✅ Fixed — added `WithImageRetry()` and `WithImageRetryDelay()` options |

### Minor

| # | Issue | Status |
|---|---|---|
| 287 | Redundant bool conversion | ✅ Fixed — removed `bool()` wrapper |
| 288 | Potential Dockerfile conflict | ✅ Fixed — `Dockerfile` takes priority over `DockerfileInline` |
| 289 | ServiceName undocumented | ✅ Verified — already present in `doc/mapping.md` line 417 |
| 290 | Inefficient fatal warning handling | ✅ Fixed — fatal check moved before `opinionated.Apply()` |
| 291 | DNSOpts warning break | ✅ Fixed — removed `break`, now warns for all unsupported DNS opts |
| 292 | Ulimits emitted twice | ✅ Fixed — removed P1 `Ulimit=` from container.go (P2 systemd limits preferred) |
| 293 | Config mode formatting | ✅ Fixed — `os.FileMode` formatted as octal `%04o` |
