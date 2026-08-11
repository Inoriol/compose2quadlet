package opinionated

import (
	"strings"

	c2qtypes "github.com/inoriol/compose2quadlet/internal/types"
)

func ApplyReferences(units []c2qtypes.QuadletUnit, cfg *c2qtypes.Config) []c2qtypes.QuadletUnit {
	if cfg.FilePrefix == "" && cfg.ProjectName == "" {
		return units
	}
	var prefix string
	if cfg.ProjectName != "" {
		prefix = cfg.FilePrefix + cfg.ProjectName + "-"
	} else {
		prefix = cfg.FilePrefix
	}
	refKeys := map[string][]string{
		"Network": {".network"},
		"Volume":  {".volume"},
		"Image":   {".image", ".build"},
	}
	unitRefKeys := map[string]string{
		"After":    ".container",
		"Requires": ".container",
		"Wants":    ".container",
		"BindsTo":  ".container",
		"PartOf":   ".container",
	}

	for ui := range units {
		for si := range units[ui].Sections {
			for di := range units[ui].Sections[si].Directives {
				d := &units[ui].Sections[si].Directives[di]

				if suffixes, ok := refKeys[d.Key]; ok {
					for vi := range d.Values {
						for _, suffix := range suffixes {
							if strings.HasSuffix(d.Values[vi], suffix) {
								d.Values[vi] = prefix + d.Values[vi]
								break
							}
						}
					}
				}

				if suffix, ok := unitRefKeys[d.Key]; ok {
					for vi := range d.Values {
						if !strings.Contains(d.Values[vi], ".") {
							d.Values[vi] = d.Values[vi] + suffix
						}
						if !strings.HasPrefix(d.Values[vi], prefix) {
							d.Values[vi] = prefix + d.Values[vi]
						}
					}
				}
			}
		}
	}
	return units
}
