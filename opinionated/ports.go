package opinionated

import (
	"strconv"
	"strings"

	c2qtypes "github.com/inoriol/compose2quadlet/internal/types"
)

func ApplyPortOffset(units []c2qtypes.QuadletUnit, cfg *c2qtypes.Config) []c2qtypes.QuadletUnit {
	if cfg.PortOffset == 0 {
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
				if d.Key != "PublishPort" {
					continue
				}
				for vi := range d.Values {
					d.Values[vi] = offsetPort(d.Values[vi], cfg.PortOffset)
				}
			}
		}
	}

	return units
}

func offsetPort(port string, offset int) string {
	parts := strings.Split(port, ":")
	last := len(parts) - 1
	targetProto := parts[last]
	proto := ""
	if idx := strings.IndexByte(targetProto, '/'); idx >= 0 {
		proto = targetProto[idx:]
		targetProto = targetProto[:idx]
	}

	if len(parts) == 1 {
		return port
	}

	publishIdx := last - 1
	publishPart := parts[publishIdx]
	if n, err := strconv.Atoi(publishPart); err == nil {
		parts[publishIdx] = strconv.Itoa(n + offset)
	}

	port = strings.Join(parts, ":")
	if proto != "" {
		port = strings.Replace(port, targetProto, targetProto+proto, 1)
	}
	return port
}
