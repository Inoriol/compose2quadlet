package compose2quadlet

import (
	"errors"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/inoriol/compose2quadlet/mapper"
	"github.com/inoriol/compose2quadlet/opinionated"
)

func Transpile(project *types.Project, opts ...TranspileOption) ([]QuadletUnit, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	var units []QuadletUnit

	for name := range project.Services {
		svc := project.Services[name]

		secretDirs := mapper.PremapSecrets(&svc, project.Secrets, project.Configs, cfg)
		project.Services[name] = svc

		var sections []Section

		unitDirs := mapper.Unit(svc)
		if len(unitDirs) > 0 {
			sections = append(sections, Section{Name: SectionUnit, Directives: unitDirs})
		}

		containerDirs := mapper.Container(svc, cfg)
		hcDirs := mapper.Healthcheck(svc, cfg)
		containerDirs = append(containerDirs, hcDirs...)
		containerDirs = append(containerDirs, secretDirs...)

		if len(containerDirs) > 0 {
			sections = append(sections, Section{Name: SectionContainer, Directives: containerDirs})
		}

		serviceDirs := mapper.Service(svc, cfg)
		depServiceDirs := mapper.UnitService(svc)
		serviceDirs = append(serviceDirs, depServiceDirs...)
		if len(serviceDirs) > 0 {
			sections = append(sections, Section{Name: SectionService, Directives: serviceDirs})
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

	imagUnits := mapper.Images(project.Services, cfg)
	units = append(units, imagUnits...)

	buildUnits := mapper.Builds(project.Services, cfg)
	units = append(units, buildUnits...)

	netUnits := mapper.Networks(project.Networks, cfg)
	units = append(units, netUnits...)

	volUnits := mapper.Volumes(project.Volumes, cfg)
	units = append(units, volUnits...)

	for _, w := range cfg.Warnings {
		if w.Level == WarningFatal {
			return nil, errors.New(w.Message)
		}
	}

	units = opinionated.Apply(units, cfg)

	return units, nil
}
