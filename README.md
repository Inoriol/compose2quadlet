# compose2quadlet

Go library that transpiles Docker Compose files into Podman Quadlet unit files.

**Not a CLI tool.** This is a library designed to be consumed by [comquad](https://github.com/inoriol/comquad) and other tools that need programmatic compose-to-quadlet conversion.

Heavily inspired by [podlet](https://github.com/containers/podlet), but designed for the Go ecosystem with structured output types instead of fragile text manipulation.

## Motivation

Replaces the `podlet` binary dependency with a native Go library. Eliminates the fragile text-based pipeline (preprocess → strip → transpile → cook → graft) with structured, type-safe quadlet output.

## Usage

```go
import (
    c2q "github.com/inoriol/compose2quadlet"
    "github.com/compose-spec/compose-go/v2/cli"
)

func main() {
    project, _ := cli.NewProjectOptions(
        []string{"compose.yaml"},
        cli.WithOsEnv,
        cli.WithDotEnv,
    ).LoadProject(context.Background())

    units, _ := c2q.Transpile(project,
        c2q.WithProjectName("myapp"),
        c2q.WithPortOffset(10000),
    )

    for _, u := range units {
        fmt.Printf("%s.%s\n", u.Name, u.Type)
    }
}
```

## Dependencies

- `github.com/compose-spec/compose-go/v2` — canonical compose parsing
- Standard library only otherwise

## License

MIT
