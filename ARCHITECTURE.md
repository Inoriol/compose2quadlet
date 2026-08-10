# Architecture Guide

This document describes the internal design, conventions, and data flow of compose2quadlet. It is written for AI agents and new contributors reading the codebase for the first time. Keep it updated as the project evolves.

## Official Specification Links

| Spec | URL |
|---|---|
| **Compose Specification** | https://github.com/compose-spec/compose-spec/blob/master/spec.md |
| **Compose Build Spec** | https://github.com/compose-spec/compose-spec/blob/master/build.md |
| **Compose Deploy Spec** | https://github.com/compose-spec/compose-spec/blob/master/deploy.md |
| **Quadlet (podman-systemd.unit.5)** | https://github.com/podman-container-tools/podman/blob/main/docs/source/markdown/podman-systemd.unit.5.md |
| **systemd.unit** | https://github.com/systemd/systemd/blob/main/man/systemd.unit.xml |
| **systemd.service** | https://github.com/systemd/systemd/blob/main/man/systemd.service.xml |
| **systemd.resource-control** | https://github.com/systemd/systemd/blob/main/man/systemd.resource-control.xml |
| **systemd.exec** | https://github.com/systemd/systemd/blob/main/man/systemd.exec.xml |
| **compose-go (Go types)** | https://github.com/compose-spec/compose-go |
| **Podman Release Notes** | https://github.com/podman-container-tools/podman/blob/main/RELEASE_NOTES.md |

## Version Tracking

Every compose→quadlet mapping is tracked with a minimum podman/systemd version in `doc/mapping.md`.
The `Since` column records when the target directive was introduced.

**Minimum baseline: Podman 4.8.0** (required for `.image` quadlet support).
Quadlet types by introduction: `.container`/`.network`/`.volume` 4.4.0, `.image` 4.8.0, `.build` 5.2.0.

When adding a new field mapping, check the [Podman Release Notes](https://github.com/podman-container-tools/podman/blob/main/RELEASE_NOTES.md)
to determine the minimum podman version. For systemd directives, the version is noted in the `Added in version X`
footer of each option in `systemd.resource-control`.

## Project Scope

**Library, not a CLI.** No `main()`, no cobra, no command-line interface. The single entry point is `Transpile()`.

**Narrow focus:** compose → quadlet only. No Kubernetes, no pods, no artifacts. Only quadlet types relevant to compose: `.container`, `.network`, `.volume`, `.image`, `.build`.

**Consumer:** [comquad](https://github.com/inoriol/comquad) imports this library. comquad keeps the CLI, orchestration (`up`/`down`/`start`/`stop`/`logs`), state management, and D-Bus communication. The library absorbs the `internal/{preprocess,transpile,cook,graft}` pipeline.

## Package Structure

```
compose2quadlet/
├── ARCHITECTURE.md           # This file
├── README.md                 # End-user documentation
├── go.mod                    # module github.com/inoriol/compose2quadlet
│
├── types.go                  # Type aliases re-exporting from internal/types/
├── transpile.go              # Entry point: Transpile(project, opts...) → ([]QuadletUnit, error)
├── options.go                # TranspileOption constructors, delegates to internal/types/
│
├── internal/
│   └── types/                # Shared types — no import cycles
│       ├── core.go           # QuadletUnit, Section, Directive, Warning, WarningLevel, UnitType, section constants
│       └── config.go         # Config, Version, Option, DefaultConfig(), Warn()
│
├── mapper/                   # Field mapping logic (implemented)
│   ├── container.go          # Container(), t0Container(), t1Container() — P1 [Container] directives
│   ├── unit.go               # Unit() — depends_on → [Unit] After=/Requires=/Wants=/BindsTo=
│   ├── healthcheck.go        # Healthcheck() — healthcheck directives
│   ├── security.go           # SecurityOpts() — security_opt parsing
│   └── ports.go              # formatPort() helper
│
├── opinionated/              # Opinionated transforms (to be implemented)
│
├── serialization/            # Serialization / deserialization (implemented)
│   └── ini.go                # Marshal(), Write(), Unmarshal()
│
├── doc/
│   └── mapping.md            # Complete field-by-field mapping reference
│
└── (reference files)
    ├── podman-systemd.unit.5.md # Full quadlet spec
    ├── systemd.unit
    ├── systemd.service
    ├── systemd.resource-control
    ├── systemd.scope
    ├── systemd.slice
    └── deploy.md
```
│   └── mapping.md            # Complete field-by-field mapping reference
│
└── (reference files)
    ├── podman-systemd.unit.5.md # Full quadlet spec
    ├── systemd.unit
    ├── systemd.service
    ├── systemd.resource-control
    ├── systemd.scope
    ├── systemd.slice
    └── deploy.md
```

## Core Types

```go
// UnitType identifies the kind of quadlet unit file.
type UnitType string  // "container" | "network" | "volume" | "image" | "build"

// QuadletUnit is a structured representation of a quadlet unit file
// before serialization to ini-format.
type QuadletUnit struct {
    Type     UnitType
    Name     string     // base name, before prefixing (e.g. "web", "db")
    Sections []Section
}

// Section represents an ini section like [Container], [Service], [Unit].
type Section struct {
    Name       string     // e.g. "Container", "Service", "Unit", "Install"
    Directives []Directive
}

// Directive is a single key-value entry. Values stores multiple
// values for directives that repeat on separate lines with the same key
// (e.g. Environment=, Volume=, PublishPort=).
type Directive struct {
    Key    string
    Values []string  // empty = key present with no value (boolean flag)
}
```

### Section Constants

Predefined section names used across the codebase:
```
SectionUnit      = "Unit"
SectionService   = "Service"
SectionInstall   = "Install"
SectionContainer = "Container"
SectionNetwork   = "Network"
SectionVolume    = "Volume"
SectionImage     = "Image"
SectionBuild     = "Build"
```

## Data Flow

```
compose.yaml
    │
    ▼
compose-go/loader.Load()  ← handles interpolation, env resolution, extension merging
    │
    ▼
*types.Project   ← canonical, fully-resolved compose model
    │
    ▼
Transpile(project, opts...)
    │
    ├── 1. Apply transpileConfig (defaults + user overrides from opts)
    │
    ├── 2. Secrets pre-mapping intercept
    │       secrets → Volume= / Secret= in [Container]
    │       configs → Mount=type=bind in [Container]
    │       (strips secrets/configs from model before field mapping)
    │
    ├── 3. Field mapping phase (mapper/)
    │       For each service:
    │       ├── container.go: service fields → [Container] directives (P1) ✅
    │       ├── unit.go: depends_on → [Unit] After=/Requires= ✅
    │       ├── healthcheck.go: healthcheck → HealthCmd= etc. ✅
    │       ├── service.go:  deploy/resources → [Service] directives (P2)
    │       ├── network.go:  networks → Network= + .network quadlet
    │       ├── volume.go:   volumes  → Volume= + .volume quadlet
    │       ├── image.go:    image    → .image quadlet
    │       └── build.go:    build    → .build quadlet
    │
    ├── 4. Top-level networks/volumes → .network/.volume quadlets
    │
    ├── 5. Opinionated transforms (opinionated/)
    │       ├── prefix.go:       cq-<project>- prefix on all unit names
    │       ├── references.go:   rewrite Network=, Volume=, After= references
    │       ├── aliases.go:      inject NetworkAlias=<service> per network
    │       ├── selinux.go:      add :z to volume mounts
    │       ├── labels.go:       inject com.comquad labels
    │       ├── network.go:      inject default network if needed
    │       ├── ports.go:        apply port offset
    │       └── install.go:      add [Install] section
    │
    ▼
[]QuadletUnit   ← structured, typed output
    │
    ▼
serialization/ini.go    ← serialize to ini text format (optional; comquad may serialize itself)
    │
    ▼
foo.container, bar.network, baz.volume files written to disk
```

## Field Mapping Priority System

Every compose field maps to exactly one of four levels. This is documented exhaustively in `doc/mapping.md`.

| P | Name | Target | Example |
|---|---|---|---|
| **1** | Direct Quadlet | `[Container]`, `[Network]`, `[Volume]`, `[Image]`, `[Build]` | `ports` → `PublishPort=` |
| **2** | Systemd | `[Service]` (resource-control, restart) or `[Unit]` (deps) | `mem_limit` → `MemoryMax=` |
| **3** | PodmanArgs | `PodmanArgs=` in `[Container]` | `tty` → `PodmanArgs=--tty` |
| **4** | Unsupported | Ignored (or warned) | `deploy.replicas`, `extends` |
| — | Structural | Generates separate quadlet unit | `build` → `.build` unit |

Priority 2 fields go in the `[Service]` section of the container unit. Quadlet passes `[Service]` directives through to the generated `.service` file, so systemd enforces them at the cgroup level. This is superior to using `PodmanArgs` because systemd can enforce limits even if podman is bypassed.

## Version Awareness

The library tracks a target podman version (via `WithPodmanVersion()` or defaults to "latest"). Mappers gate output based on this version:

| Target version | `entrypoint:` behavior | `build:` behavior |
|---|---|---|
| Zero / latest | Emit `Entrypoint=` (P1) | Emit `.build` unit (P1) |
| `5.0.0` | Emit `Entrypoint=` (P1, since 5.0) | Emit `.build` unit (P1, since 5.2) |
| `4.8.0` | Emit `PodmanArgs=--entrypoint ...` (P3 fallback) | **Fatal error** — impossible |

The version check is centralized:

- **`mapper/mapper.go`** — A central version-gated field registry handles simple 1:1 fields where the P1 key is available from a minimum version and the P3 `PodmanArgs` fallback is mechanical. ~85% of version-gated fields live here.
- **Custom mapper functions** — Complex fields (e.g. `build` which produces a whole unit, `entrypoint` which has a non-trivial P3 format) use `cfg.podmanVersion.AtLeast(...)` directly in their code.

Systemd versions are tracked only as documentation in `doc/mapping.md`. In practice, modern podman implies modern systemd, so the library collapses to a single podman version axis.

## Warning System

Every field that cannot be mapped is surfaced — there are **no silent skips**. Three severity levels:

| Level | Meaning | Example | Consumer impact |
|---|---|---|---|
| `WarningSkipped` | Feature unavailable at target podman version | `network_aliases` on podman 4.8.0 | Info: "field skipped, requires podman 5.2.0" |
| `WarningDegraded` | P3 PodmanArgs fallback instead of P1 | `entrypoint` on podman 4.8.0 | Warn: "mapped via PodmanArgs, upgrade to podman 5.0 for native support" |
| `WarningFatal` | Mapping is impossible at this version | `build:` on podman 4.8.0 | `Transpile()` returns error |

Warnings are collected in `transpileConfig.Warnings` during the pipeline and surfaced alongside the result. Consumers (comquad) decide how to present each level. No separate error/warning channel is threaded through mapper signatures — the config acts as a shared collector.

```go
// Internal usage in a mapper:
cfg.warn(Warning{
    Level:   WarningDegraded,
    Service:  svc.Name,
    Field:    "entrypoint",
    Message:  "using PodmanArgs fallback",
    Since:    "5.0.0",
})
```

## Opinionated Defaults

All transforms are **enabled by default** and can be individually disabled via `TranspileOption`. This matches comquad's current behavior.

| Transform | Option to disable | What it does |
|---|---|---|
| File prefixing | `WithoutPrefix()` | Prepends `cq-<project>-` to unit filenames |
| Reference rewriting | *(always on)* | Rewrites `Network=`, `Volume=`, `After=`, `Requires=` to prefixed names |
| NetworkAlias injection | `WithoutNetworkAliases()` | Adds `NetworkAlias=<service>` for each network a service connects to |
| SELinux `:z` | `WithoutSELinux()` | Appends `:z` to all bind-mount volumes |
| Managed label | `WithoutManagedLabel()` | Adds `Label=com.comquad.managed=true` |
| Project label | *(requires `WithProjectName()`)* | Adds `Label=com.comquad.project=<name>` |
| Default network | `WithoutDefaultNetwork()` | Injects `cq-default.network` if no networks defined |
| Port offset | `WithPortOffset(N)` | Adds offset to all host-side published ports |
| AutoUpdate | `WithAutoUpdate()` | Adds `AutoUpdate=registry` to containers |
| Install section | `WithoutInstallSection()` | Adds `[Install] WantedBy=default.target` |

## Quadlet-Specific Behaviors

### Image Quadlet Generation

Every service with `image:` gets a companion `.image` quadlet. The container unit references it via `Image=<name>.image`. This splits image pulling into a separate systemd unit, enabling dependency ordering.

### Quadlet Cross-Reference Syntax

Quadlet uses special extension syntax for unit references:
- `Network=<name>.network` — references a `.network` quadlet
- `Volume=<name>.volume` — references a `.volume` quadlet
- `Image=<name>.image` — references a `.image` quadlet
- `Image=<name>.build` — references a `.build` quadlet

The mapper must output these `.network`/`.volume`/`.image`/`.build` suffixes so that systemd dependency chains are created automatically.

### Dependency Translation

`[Unit]` directives (`After=`, `Requires=`, `Wants=`, `BindsTo=`, `PartOf=`) between quadlet units are automatically translated by the quadlet generator. For example, `After=db.container` in a `web.container` unit creates a proper systemd `After=db.service` dependency.

## Conventions

### Go Code Style
- Standard library only (plus compose-go/v2). No external process execution, no additional YAML libraries.
- Package name: `compose2quadlet` (not `main` — this is a library).
- No code comments in implementation files unless the logic is non-obvious. Mapping is self-documenting via directive names.
- Tests live in `*_test.go` files alongside the code they test.

### Naming
- `QuadletUnit.Name` is the **base name** before prefixing (e.g. `"web"`, `"db"`, `"default"`).
- The full filename is constructed as `<prefix><Name>.<type>` during serialization.
- Section names are PascalCase strings matching ini section headers: `"Container"`, `"Service"`, `"Unit"`.
- Directive keys are the exact quadlet/systemd directive names: `"PublishPort="`, `"MemoryMax="`.

### Directive Ordering
Directives within a section should follow the order from the quadlet spec where possible. Systemd `[Service]` directives go after all `[Container]` directives. `[Unit]` goes before `[Container]`, `[Install]` goes last.

## Testing Strategy

### Tier 0 — Serialization (`serialization/ini.go`)
QuadletUnit → ini text correctness:

- Section ordering (`[Unit]` → `[Container]` → `[Service]` → `[Install]`)
- Multi-value directive rendering (multiple `Volume=`, `Environment=` lines)
- Empty-value directives (boolean flags: `NoNewPrivileges=`)
- **Empty-default**: each unit type with only mandatory fields renders correctly (no trailing empty lines, no missing `\n`)
- **Empty section omission**: sections with zero directives are not rendered (no dangling `[Unit]\n` header)
- **Round-trip**: serialize → deserialize → identical `QuadletUnit`, for every unit type
- **Comment line**: `# FileName=<name>` header line present/absent in output

### Tier 1 — Mapper Unit Tests (`mapper/*_test.go`)
Each mapper function tested independently with table-driven tests.
Input: compose-go `types.ServiceConfig` (or relevant field subset).
Output: expected `[]Directive` (or `[]QuadletUnit` for structural mappers).

Tests focus on **correct output at latest podman version**. Version-gated behavior is tested in Tier 2.

### Tier 2 — Version Matrix (`version_test.go` in package root)
A dedicated test file that runs the same compose input through multiple target podman versions and asserts correct behavior per version boundary (4.8.0, 5.0.0, 5.2.0, 5.3.0, 5.5.0, 5.8.0, 6.0.0, 6.1.0). Covers:

- Fields that promote from P3 PodmanArgs to P1 native directive (`entrypoint`, `stop_signal`, `extra_hosts`)
- Fields that become available at a boundary (`notify`, `network_aliases`, `addhost`, `log_options`)
- Fields that switch section (`mem_limit`: `[Service] MemoryMax=` → `[Container] Memory=`)
- Structural blocks that cause fatal errors (`build` on < 5.2.0)
- Correct `Warning` collection at each severity level

### Tier 3 — Pipeline Integration (`transpile_test.go`)
Full compose YAML → compose-go parse → `Transpile()` → verify `[]QuadletUnit` structure:
- Single service, multi-service, no-service (top-level networks/volumes only)
- Option combinatorics: pairs of enable/disable on opinionated transforms
- Warning collection verification
- Edge cases: empty service, build-only service, host network mode, external volumes/networks

```go
func TestTranspile_SimpleWeb(t *testing.T) {
    project := loadFixture(t, "testdata/simple-web.yaml")
    units, err := Transpile(project, WithProjectName("test"))
    require.NoError(t, err)
    require.Len(t, units, 2) // web.container + web.image
    assertSectionKey(t, units[0], SectionContainer, "PublishPort")
}
```

### Tier 4 — End-to-End (deferred to comquad)
comquad's existing `tests/integration/` harness. The library itself does not start podman or systemd.

### Test Conventions
- Fixture compose files live in `testdata/` at the package root.
- Table-driven tests use `t.Run()` for each entry with descriptive names.
- Golden files for serialization live in `testdata/serialization/` with `.golden` extension.
- Test helper functions (`assertSectionKey`, `loadFixture`) are shared in a `helpers_test.go` file.
- No external test dependencies beyond the standard library and compose-go/v2.
- **Empty-default pattern**: every unit type gets a test verifying that a `QuadletUnit` with only mandatory fields serializes correctly — catches section rendering bugs early.
- **Round-trip pattern**: for serialization, every test verifying serialization should also verify deserialization produces the same struct.

## Development Order (Milestones)

From the project scope document:

1. **MVP** — `.container` files only, priority-1 field mappings, no opinionated transforms ✅
2. **Full compose parity** — `.network`, `.volume`, `.image`, `.build` support, all priority-1 + priority-2
3. **Opinionated defaults** — all comquad transforms ported as opt-out `TranspileOption`s
4. **Deploy + systemd** — `deploy.resources`, `deploy.restart_policy` mapped to `[Service]`
5. **Secrets + builds** — compose `secrets:` and `build:` handled natively
6. **Integration** — comquad imports the library, drops podlet dependency
7. **Deprecate podlet** — comquad no longer requires podlet binary at runtime

## Key Design Decisions

### Why structured output instead of text?
Podlet emits ini text. Manipulating text is fragile (regex, line parsing). `QuadletUnit` structs allow programmatic modification — rename files, rewrite references, inject labels, offset ports — before serialization. This eliminates comquad's entire strip/cook/graft pipeline.

### Why not use quadlet directives for everything?
Some compose fields have better systemd equivalents. For example, `mem_limit` could be `PodmanArgs=--memory ...` (P3) but `MemoryMax=` in `[Service]` (P2) is enforced at the cgroup level by systemd itself. Priority 2 always wins over priority 3 when both are possible.

### Why separate `.image` quadlets?
Splitting image pulls into a separate unit enables proper dependency ordering. The `.image` unit completes before the `.container` unit starts. This also enables `AutoUpdate=registry` on the image unit, triggering updates independently.

### Why pre-mapping intercept for secrets?
Secrets and configs need different handling depending on their type (external vs file vs environment). Intercepting them before the field mapper runs avoids collision with the volume mapper and ensures correct `Secret=` vs `Volume=` routing.
