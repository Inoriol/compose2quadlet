package types

type Config struct {
	SelinuxContext  bool
	FilePrefix      string
	DefaultNetwork  bool
	PortOffset      int
	ProjectName     string
	Labels          map[string]string
	AutoUpdate      bool
	InstallSection  bool
	NetworkAliases  bool
	PodmanVersion   Version
	Warnings        []Warning
}

type Version struct {
	Major, Minor, Patch int
}

func (v Version) AtLeast(major, minor int) bool {
	if v.Major == 0 && v.Minor == 0 {
		return true
	}
	if v.Major > major {
		return true
	}
	if v.Major == major && v.Minor >= minor {
		return true
	}
	return false
}

func DefaultConfig() *Config {
	return &Config{
		SelinuxContext: true,
		FilePrefix:     "cq-",
		DefaultNetwork: true,
		PortOffset:     0,
		AutoUpdate:     false,
		InstallSection: true,
		NetworkAliases: true,
	}
}

type Option func(*Config)

func WithoutSELinux() Option {
	return func(c *Config) { c.SelinuxContext = false }
}

func WithoutPrefix() Option {
	return func(c *Config) { c.FilePrefix = "" }
}

func WithPrefix(prefix string) Option {
	return func(c *Config) { c.FilePrefix = prefix }
}

func WithoutDefaultNetwork() Option {
	return func(c *Config) { c.DefaultNetwork = false }
}

func WithPortOffset(offset int) Option {
	return func(c *Config) { c.PortOffset = offset }
}

func WithProjectName(name string) Option {
	return func(c *Config) { c.ProjectName = name }
}

func WithLabels(labels map[string]string) Option {
	return func(c *Config) { c.Labels = labels }
}

func WithAutoUpdate() Option {
	return func(c *Config) { c.AutoUpdate = true }
}

func WithoutInstallSection() Option {
	return func(c *Config) { c.InstallSection = false }
}

func WithoutNetworkAliases() Option {
	return func(c *Config) { c.NetworkAliases = false }
}

func WithPodmanVersion(v Version) Option {
	return func(c *Config) { c.PodmanVersion = v }
}

func (c *Config) Warn(w Warning) {
	c.Warnings = append(c.Warnings, w)
}
