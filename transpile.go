package compose2quadlet

import "github.com/compose-spec/compose-go/v2/types"

// Transpile converts a parsed compose project into quadlet units.
// All opinionated transforms are applied by default; disable them via options.
func Transpile(project *types.Project, opts ...TranspileOption) ([]QuadletUnit, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	// TODO: implement field mapping pipeline
	_ = project
	_ = cfg
	return nil, nil
}
