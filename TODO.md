# Implementation TODO

## Code Quality Issues (from code review)

### Critical

| # | Issue | Location | Description |
|---|---|---|---|
| 279 | Non-deterministic map iteration | `mapper/container.go:176,273,276,423,453,470,500,522` + others | Maps like `Labels`, `Annotations`, `Environment`, `Sysctls`, `Logging.Options`, `LogOpt`, `Ulimits`, `ExtraHosts` produce different directive ordering on each run. Breaks reproducibility and golden file testing. **Fix**: Sort keys before iteration. |
| 280 | Build + Image logic bug | `mapper/container.go:24-42` | When `svc.Build != nil` and podman < 5.2, fatal warning emitted but no `Image=` directive generated (else-if branches skipped). Results in containers with no Image directive. **Fix**: Add fallback logic or early return after fatal warning. |
| 281 | Conflicting CPU directives | `mapper/service.go:56-69` | Both `svc.CPUs` and `svc.CPUQuota` emit `CPUQuota=` directives. If both set, they conflict and last one wins. **Fix**: Prioritize one over the other or merge correctly. |
| 282 | Conflicting restart directives | `mapper/service.go:142-153` | Both `svc.Restart` and `svc.Deploy.RestartPolicy` can emit `Restart=` directives with no precedence logic. **Fix**: Check if `Deploy.RestartPolicy` exists and skip `svc.Restart` if so. |

### Medium

| # | Issue | Location | Description |
|---|---|---|---|
| 283 | Incomplete port offset logic | `opinionated/ports.go:38-63` | Returns early for single-port format like `"80"` but doesn't handle `"80/udp"`. Complex string manipulation is error-prone. **Fix**: Parse port format more robustly. |
| 284 | SELinux false positives | `opinionated/selinux.go:27` | Uses `strings.Contains(d.Values[vi], ":z")` which matches substrings. Path like `/data:zoo:/mnt` would incorrectly skip SELinux labeling. **Fix**: Check for `:z` or `:Z` only at end or after proper parsing. |
| 285 | Missing opinionated tests | `opinionated/` package | No test files for critical transforms (prefix, references, SELinux, port offset, etc.). **Fix**: Add comprehensive tests for all opinionated transforms. |
| 286 | Hardcoded retry values | `mapper/image.go:50-53` | `Retry=3` and `RetryDelay=5s` hardcoded. TODO.md mentions these should come from config, but no config option exists. **Fix**: Add config options or document as intentional defaults. |

### Minor

| # | Issue | Location | Description |
|---|---|---|---|
| 287 | Redundant bool conversion | `mapper/network.go:13`, `mapper/volume.go:13` | `bool(nc.External)` unnecessary since `External` is already bool type. |
| 288 | Potential Dockerfile conflict | `mapper/build.go:28-30` | Both `Dockerfile` and `DockerfileInline` can emit `File=` directives. If both set, two `File=` directives emitted. |
| 289 | ServiceName undocumented | `mapper/container.go:46-58` | Emits `ServiceName=` directive but not listed in TODO.md. May be intentional but lacks documentation. |
| 290 | Inefficient fatal warning handling | `transpile.go:73-77` | Fatal warnings checked after all processing completes. Wastes CPU cycles on work that will be discarded. |
| 291 | DNSOpts warning break | `mapper/container.go:219-232` | Loop breaks after first warning, only one warning emitted even with multiple options. Inconsistent with other similar loops. |
| 292 | Ulimits emitted twice | `mapper/container.go:500-518` + `mapper/service.go:208-224` | Ulimits emit both `Ulimit=` (P1) and `LimitXXX=` (P2) directives. May cause conflicts or confusion. |
| 293 | Config mode formatting | `mapper/secrets.go:112` | `ref.Mode.String()` may not format octal modes correctly for quadlet. |
