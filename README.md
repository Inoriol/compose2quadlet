# compose2quadlet

> **ARCHIVED** — This library has been merged into the main [comquad](https://github.com/Inoriol/comquad) repository.
> The code now lives at [`comquad/compose2quadlet/`](https://github.com/Inoriol/comquad/blob/main/compose2quadlet/README.md).
> This repository is read-only and no longer maintained.

Go library that transpiles Docker Compose files into Podman Quadlet unit files.

**Not a CLI tool.** This is a library designed to be consumed by [comquad](https://github.com/Inoriol/comquad) and other tools that need programmatic compose-to-quadlet conversion.

Heavily inspired by [podlet](https://github.com/containers/podlet), but designed for the Go ecosystem.

## Motivation

Replaces the `podlet` binary dependency with a native Go library. Eliminates the fragile text-based pipeline (preprocess → strip → transpile → cook → graft) for comquad (that currently uses podlet) with structured, type-safe quadlet output.

## Usage

**Quick path** — `TranspileFile` loads the compose file with `.env` resolution and transpiles in one call:

```go
import (
    c2q "github.com/inoriol/compose2quadlet"
    "github.com/inoriol/compose2quadlet/serialization"
)

func main() {
    units, _ := c2q.TranspileFile("compose.yaml",
        c2q.WithProjectName("myapp"),
        c2q.WithPortOffset(10000),
        c2q.WithAutoUpdate(),
        c2q.WithLabels(map[string]string{
            "com.myorg.managed": "true",
            "com.myorg.project": "myapp",
        }),
    )

    // Write all units to a directory
    serialization.WriteUnits("/etc/containers/systemd", units)
}
```

**Full control** — load the compose-go project yourself and call `Transpile`:

```go
import (
    c2q "github.com/inoriol/compose2quadlet"
    "github.com/inoriol/compose2quadlet/serialization"
    "github.com/compose-spec/compose-go/v2/cli"
)

func main() {
    opts, _ := cli.NewProjectOptions(
        []string{"compose.yaml"},
        cli.WithOsEnv,
        cli.WithDotEnv,
    )
    project, _ := opts.LoadProject(context.Background())

    units, _ := c2q.Transpile(project,
        c2q.WithProjectName("myapp"),
        c2q.WithPortOffset(10000),
        c2q.WithAutoUpdate(),
        c2q.WithLabels(map[string]string{
            "com.myorg.managed": "true",
            "com.myorg.project": "myapp",
        }),
    )

    for _, u := range units {
        fmt.Printf("%s.%s\n", u.Name, u.Type)
    }

    // Serialize to ini format for writing to disk
    text := serialization.Marshal(units[0])
    fmt.Println(text)
}
```

## Packages

- **Root** — `Transpile()`, `TranspileFile()`, `ParseVersion()`, core types (`QuadletUnit`, `Section`, `Directive`), option constructors
- **`mapper/`** — Field mapping: compose-go `ServiceConfig` / `*Project` → quadlet directives and units
  - `container.go`, `healthcheck.go`, `security.go`, `ports.go` — per-service container directives
  - `service.go` — systemd `[Service]` resource-control and restart directives
  - `unit.go` — `depends_on` → `[Unit]` dependencies and health polling hooks
  - `image.go`, `build.go` — companion `.image` and `.build` structural units
  - `network.go`, `volume.go` — top-level `.network` and `.volume` structural units
  - `secrets.go` — pre-mapping interceptor for secrets and configs
  - `dockerfile.go` — `PatchDockerfileFROM()` normalizes bare image names in Dockerfile FROM lines
- **`opinionated/`** — Composable post-processing transforms (prefix, references, SELinux, labels, default network, port offset, auto-update, install section)
- **`serialization/`** — `Marshal()`, `Write()`, `WriteUnits()`, `Unmarshal()` for ini-format serialization

## Dependencies

- `github.com/compose-spec/compose-go/v2` — canonical compose parsing
- Standard library only otherwise

## License

MIT
