package opinionated

import (
	"strings"

	c2qtypes "github.com/inoriol/compose2quadlet/internal/types"
)

func ApplySELinux(units []c2qtypes.QuadletUnit, cfg *c2qtypes.Config) []c2qtypes.QuadletUnit {
	if !cfg.SelinuxContext {
		return units
	}

	for ui := range units {
		if units[ui].Type != c2qtypes.UnitContainer {
			continue
		}
		for si := range units[ui].Sections {
			if units[ui].Sections[si].Name != c2qtypes.SectionContainer {
				continue
			}
			for di := range units[ui].Sections[si].Directives {
				d := &units[ui].Sections[si].Directives[di]

				if d.Key == "Volume" || d.Key == "Mount" {
					for vi := range d.Values {
						if strings.Contains(d.Values[vi], ":z") || strings.Contains(d.Values[vi], ":Z") {
							continue
						}
						if d.Key == "Mount" {
							d.Values[vi] = d.Values[vi] + ",selinux=z"
						} else {
							d.Values[vi] = d.Values[vi] + ",z"
						}
					}
				}
			}
		}
	}

	return units
}
