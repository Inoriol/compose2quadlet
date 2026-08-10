package compose2quadlet

import (
	"errors"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/inoriol/compose2quadlet/mapper"
)

func Transpile(project *types.Project, opts ...TranspileOption) ([]QuadletUnit, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	var units []QuadletUnit

	for name, svc := range project.Services {
		var sections []Section

		unitDirs := mapper.Unit(svc)
		if len(unitDirs) > 0 {
			sections = append(sections, Section{Name: SectionUnit, Directives: unitDirs})
		}

		containerDirs := mapper.Container(svc, cfg)
		hcDirs := mapper.Healthcheck(svc, cfg)
		containerDirs = append(containerDirs, hcDirs...)

		if len(containerDirs) > 0 {
			sections = append(sections, Section{Name: SectionContainer, Directives: containerDirs})
		}

		if len(sections) == 0 {
			continue
		}

		units = append(units, QuadletUnit{
			Type:     UnitContainer,
			Name:     name,
			Sections: sections,
		})
	}

	for _, w := range cfg.Warnings {
		if w.Level == WarningFatal {
			return nil, errors.New(w.Message)
		}
	}

	return units, nil
}
