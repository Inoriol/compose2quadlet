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
├── types.go                  # Core types: QuadletUnit, Section, Directive, UnitType
├── transpile.go              # Entry point: Transpile(project, opts...) → []QuadletUnit
├── options.go                # transpileConfig, TranspileOption functions
│
├── mapper/                   # Field mapping logic (to be implemented)
│   ├── container.go          # service → [Container] directives (priority 1)
│   ├── service.go            # deploy/restart/resource → [Service] directives (priority 2)
│   ├── unit.go               # depends_on → [Unit] After=/Requires=/Wants=/BindsTo=
│   ├── network.go            # service networks → [Container] Network= + .network quadlet
│   ├── volume.go             # service volumes → [Container] Volume=/Mount= + .volume quadlet
│   ├── image.go              # image/platform → .image quadlet
│   ├── build.go              # build blocks → .build quadlet
│   └── secrets.go            # secrets/configs → Volume= / Secret= (pre-mapping intercept)
│
├── opinionated/              # Opinionated transforms (to be implemented)
│   ├── prefix.go             # cq-<project>- file prefixing
│   ├── references.go         # Cross-unit reference rewriting
│   ├── aliases.go            # NetworkAlias=<service> per network
│   ├── selinux.go            # :z on volume mounts
│   ├── labels.go             # com.comquad.managed / com.comquad.project labels
│   ├── network.go            # Default network injection (cq-default.network)
│   ├── ports.go              # Rootless port offsetting
│   └── install.go            # [Install] WantedBy=default.target
│
├── serde/                    # Serialization / deserialization (to be implemented)
│   └── ini.go                # QuadletUnit → .container/.volume/.network file text
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
    │       ├── container.go: service fields → [Container] directives (P1)
    │       ├── service.go:  deploy/resources → [Service] directives (P2)
    │       ├── unit.go:     depends_on → [Unit] After=/Requires= (P1/P2)
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
serde/ini.go    ← serialize to ini text format (optional; comquad may serialize itself)
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

### Unit Tests
Each mapper function is tested independently with table-driven tests. Input: a `types.ServiceConfig` (or subset). Output: expected `[]Directive`.

### Compatibility Parity
The podlet compatibility matrix (`podlet-v0.3.2-podman-v5.8.json`, 987 test entries) serves as a baseline. For every field where podlet achieves `level 2` (Supported), the library must achieve at least the same level.

### Integration Tests
Deferred to comquad's existing `tests/integration/` harness. The library itself doesn't start podman or systemd.

## Development Order (Milestones)

From the project scope document:

1. **MVP** — `.container` files only, priority-1 field mappings, no opinionated transforms
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
