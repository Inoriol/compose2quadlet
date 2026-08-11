package mapper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/types"
	c2qtypes "github.com/inoriol/compose2quadlet/internal/types"
)

func Builds(services types.Services, cfg *c2qtypes.Config) []c2qtypes.QuadletUnit {
	var units []c2qtypes.QuadletUnit
	for name, svc := range services {
		if svc.Build == nil {
			continue
		}
		if !cfg.PodmanVersion.AtLeast(5, 2) {
			continue
		}
		build := svc.Build
		var dirs []c2qtypes.Directive

		if build.Context != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "SetWorkingDirectory", Values: []string{build.Context}})
		}

		if cfg.NormalizeDockerfile && build.Dockerfile != "" && cfg.BuildCacheDir != "" {
			content, err := os.ReadFile(build.Dockerfile)
			if err == nil {
				patched, err := PatchDockerfileFROM(bytes.NewReader(content))
				if err == nil {
					patchedPath := filepath.Join(cfg.BuildCacheDir, name+".Dockerfile")
					os.MkdirAll(filepath.Dir(patchedPath), 0755)
					os.WriteFile(patchedPath, patched, 0644)
					build.Dockerfile = patchedPath
				}
			}
		}

		if build.Dockerfile != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "File", Values: []string{build.Dockerfile}})
		} else if build.DockerfileInline != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "File", Values: []string{build.DockerfileInline}})
		}
		if build.Target != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "Target", Values: []string{build.Target}})
		}
		if build.Network != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "Network", Values: []string{build.Network}})
		}
		if build.NoCache {
			dirs = append(dirs, c2qtypes.Directive{Key: "PodmanArgs", Values: []string{"--no-cache"}})
		}
		for _, k := range sortedKeys(build.Labels) {
			dirs = append(dirs, c2qtypes.Directive{Key: "Label", Values: []string{fmt.Sprintf("%s=%s", k, build.Labels[k])}})
		}
		for _, tag := range build.Tags {
			dirs = append(dirs, c2qtypes.Directive{Key: "ImageTag", Values: []string{tag}})
		}
		for _, secret := range build.Secrets {
			if secret.Source != "" {
				dirs = append(dirs, c2qtypes.Directive{Key: "Secret", Values: []string{secret.Source}})
			}
		}
		if len(build.Args) > 0 {
			if cfg.PodmanVersion.AtLeast(5, 7) {
				for _, k := range sortedKeys(build.Args) {
					v := build.Args[k]
					if v != nil {
						dirs = append(dirs, c2qtypes.Directive{Key: "BuildArg", Values: []string{k + "=" + *v}})
					} else {
						dirs = append(dirs, c2qtypes.Directive{Key: "BuildArg", Values: []string{k}})
					}
				}
			} else {
				cfg.Warn(c2qtypes.Warning{
					Level:   c2qtypes.WarningSkipped,
					Service: name,
					Field:   "build.args",
					Message: "requires podman >= 5.7.0",
					Since:   "5.7.0",
				})
			}
		}

		if len(build.Platforms) > 0 {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: name,
				Field:   "build.platforms",
				Message: "multi-arch build not applicable",
			})
		}
		if len(build.ExtraHosts) > 0 {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: name,
				Field:   "build.extra_hosts",
				Message: "no podman build equivalent",
			})
		}

		units = append(units, c2qtypes.QuadletUnit{
			Type:     c2qtypes.UnitBuild,
			Name:     name,
			Sections: []c2qtypes.Section{{Name: c2qtypes.SectionBuild, Directives: dirs}},
		})
	}
	return units
}
