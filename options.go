package compose2quadlet

import "github.com/inoriol/compose2quadlet/internal/types"

func defaultConfig() *types.Config {
	return types.DefaultConfig()
}

func WithoutSELinux() types.Option     { return types.WithoutSELinux() }
func WithoutPrefix() types.Option       { return types.WithoutPrefix() }
func WithPrefix(p string) types.Option  { return types.WithPrefix(p) }
func WithoutDefaultNetwork() types.Option   { return types.WithoutDefaultNetwork() }
func WithPortOffset(o int) types.Option     { return types.WithPortOffset(o) }
func WithProjectName(n string) types.Option { return types.WithProjectName(n) }
func WithLabels(l map[string]string) types.Option { return types.WithLabels(l) }
func WithAutoUpdate() types.Option          { return types.WithAutoUpdate() }
func WithoutInstallSection() types.Option   { return types.WithoutInstallSection() }
func WithoutNetworkAliases() types.Option   { return types.WithoutNetworkAliases() }
func WithPodmanVersion(v Version) types.Option { return types.WithPodmanVersion(v) }
