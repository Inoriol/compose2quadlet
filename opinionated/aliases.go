package opinionated

import c2qtypes "github.com/inoriol/compose2quadlet/internal/types"

func ApplyNetworkAliases(units []c2qtypes.QuadletUnit, cfg *c2qtypes.Config) []c2qtypes.QuadletUnit {
	if !cfg.NetworkAliases {
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

			var networks []string
			for _, d := range units[ui].Sections[si].Directives {
				if d.Key == "Network" {
					for _, v := range d.Values {
						networks = append(networks, v)
					}
				}
			}

			seen := map[string]bool{}
			for _, net := range networks {
				name := stripSuffix(net)
				if !seen[name] {
					seen[name] = true
					units[ui].Sections[si].Directives = append(units[ui].Sections[si].Directives,
						c2qtypes.Directive{Key: "NetworkAlias", Values: []string{name}})
				}
			}
		}
	}

	return units
}

func stripSuffix(s string) string {
	for _, suffix := range []string{".network", ".volume", ".image", ".build", ".container"} {
		if len(s) > len(suffix) && s[len(s)-len(suffix):] == suffix {
			return s[:len(s)-len(suffix)]
		}
	}
	return s
}
