package mapper

import (
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
	c2qtypes "github.com/inoriol/compose2quadlet/internal/types"
)

func TestContainer_Image(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Image: "nginx:latest"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Image", "nginx:latest")
}

func TestContainer_ContainerName(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", ContainerName: "my-web"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "ContainerName", "my-web")
}

func TestContainer_Init(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Init: boolPtr(true)}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "RunInit", "true")
}

func TestContainer_InitNil(t *testing.T) {
	svc := types.ServiceConfig{Name: "s", Init: nil}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	if hasDirective(dirs, "RunInit", "true") {
		t.Fatal("should not emit RunInit when Init is nil")
	}
}

func TestContainer_InitFalse(t *testing.T) {
	svc := types.ServiceConfig{Name: "s", Init: boolPtr(false)}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	if hasDirective(dirs, "RunInit", "true") {
		t.Fatal("should not emit RunInit when Init is false")
	}
}

func TestContainer_ReadOnly(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", ReadOnly: true}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "ReadOnly", "true")
}

func TestContainer_User(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", User: "1000"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "User", "1000")
}

func TestContainer_Labels(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Labels: types.Labels{"env": "prod", "team": "infra"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Label", "env=prod")
	assertDirective(t, dirs, "Label", "team=infra")
}

func TestContainer_Annotations(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Annotations: types.Mapping{"foo": "bar"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Annotation", "foo=bar")
}

func TestContainer_CapAdd(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", CapAdd: []string{"NET_ADMIN", "SYS_PTRACE"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "AddCapability", "NET_ADMIN")
	assertDirective(t, dirs, "AddCapability", "SYS_PTRACE")
}

func TestContainer_CapDrop(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", CapDrop: []string{"NET_RAW"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "DropCapability", "NET_RAW")
}

func TestContainer_Expose(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Expose: types.StringOrNumberList{"80", "443"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "ExposeHostPort", "80")
	assertDirective(t, dirs, "ExposeHostPort", "443")
}

func TestContainer_GroupAdd(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", GroupAdd: []string{"999"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "GroupAdd", "999")
}

func TestContainer_GroupAdd_Unavailable(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", GroupAdd: []string{"999"}}
	cfg := c2qtypes.DefaultConfig()
	cfg.PodmanVersion = c2qtypes.Version{Major: 4, Minor: 8}
	dirs := Container(svc, cfg)

	if len(dirs) > 0 {
		t.Fatal("expected no directives for group_add at 4.8")
	}
	if !hasWarning(cfg, "web", "group_add") {
		t.Fatal("expected WarningSkipped for group_add")
	}
}

func TestContainer_WorkingDir(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", WorkingDir: "/app"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "WorkingDir", "/app")
}

func TestContainer_ServiceName(t *testing.T) {
	svc := types.ServiceConfig{Name: "my-web"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "ServiceName", "my-web")
}

func TestContainer_SecurityOpt_NoNewPrivileges(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"no-new-privileges"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "NoNewPrivileges", "")
}

func TestContainer_SecurityOpt_LabelDisable(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"label=disable"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "SecurityLabelDisable", "")
}

func TestContainer_SecurityOpt_LabelNested(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"label=nested"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "SecurityLabelNested", "")
}

func TestContainer_SecurityOpt_Seccomp(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"seccomp=/etc/seccomp.json"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "SeccompProfile", "/etc/seccomp.json")
}

func TestContainer_SecurityOpt_LabelType(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"label=type:spc_t"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "SecurityLabelType", "spc_t")
}

func TestContainer_SecurityOpt_LabelLevel(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"label=level:s0:c1,c2"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "SecurityLabelLevel", "s0:c1,c2")
}

func TestContainer_SecurityOpt_LabelFiletype(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"label=filetype:container_t"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "SecurityLabelFileType", "container_t")
}

func TestContainer_SecurityOpt_Mask(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"mask=/proc/keys"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Mask", "/proc/keys")
}

func TestContainer_SecurityOpt_Unmask(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"unmask=all"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Unmask", "all")
}

func TestContainer_SecurityOpt_AppArmor(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"apparmor=my-profile"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "AppArmor", "my-profile")
}

func TestContainer_AppArmor_Unavailable(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", SecurityOpt: []string{"apparmor=my-profile"}}
	cfg := c2qtypes.DefaultConfig()
	cfg.PodmanVersion = c2qtypes.Version{Major: 5, Minor: 5}
	dirs := Container(svc, cfg)

	if hasDirective(dirs, "AppArmor", "my-profile") {
		t.Fatal("expected no AppArmor directive at 5.5")
	}
	if !hasWarning(cfg, "web", "security_opt") {
		t.Fatal("expected WarningSkipped for apparmor")
	}
}

func TestContainer_UserNS(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", UserNSMode: "keep-id"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "UserNS", "keep-id")
}

func TestContainer_Hostname(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Hostname: "myhost"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "HostName", "myhost")
}

func TestContainer_ShmSize(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", ShmSize: 67108864}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "ShmSize", "67108864")
}

func TestContainer_Sysctls(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Sysctls: types.Mapping{"net.core.somaxconn": "1024", "vm.overcommit": "1"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Sysctl", "net.core.somaxconn=1024")
	assertDirective(t, dirs, "Sysctl", "vm.overcommit=1")
}

func TestContainer_DNS(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DNS: types.StringList{"8.8.8.8", "1.1.1.1"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "DNS", "8.8.8.8")
	assertDirective(t, dirs, "DNS", "1.1.1.1")
}

func TestContainer_DNSSearch(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DNSSearch: types.StringList{"example.com", "internal"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "DNSSearch", "example.com")
	assertDirective(t, dirs, "DNSSearch", "internal")
}

func TestContainer_DNSOpt(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DNSOpts: []string{"ndots:2", "timeout:1"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "DNSOption", "ndots:2")
	assertDirective(t, dirs, "DNSOption", "timeout:1")
}

func TestContainer_StopGracePeriod(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", StopGracePeriod: durationPtr(types.Duration(30_000_000_000))}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "StopTimeout", "30s")
}

func TestContainer_StopSignal(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", StopSignal: "SIGQUIT"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "StopSignal", "SIGQUIT")
}

func TestContainer_PullPolicy(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", PullPolicy: "always"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Pull", "always")
}

func TestContainer_VersionGatedWarnings(t *testing.T) {
	tests := []struct {
		name    string
		svc     types.ServiceConfig
		version c2qtypes.Version
		field   string
	}{
		{"working_dir <4.6", types.ServiceConfig{Name: "s", WorkingDir: "/app"}, c2qtypes.Version{Major: 4, Minor: 5}, "working_dir"},
		{"userns_mode <4.5", types.ServiceConfig{Name: "s", UserNSMode: "keep-id"}, c2qtypes.Version{Major: 4, Minor: 4}, "userns_mode"},
		{"hostname <4.6", types.ServiceConfig{Name: "s", Hostname: "h"}, c2qtypes.Version{Major: 4, Minor: 5}, "hostname"},
		{"shm_size <4.7", types.ServiceConfig{Name: "s", ShmSize: 1024}, c2qtypes.Version{Major: 4, Minor: 6}, "shm_size"},
		{"sysctls <4.6", types.ServiceConfig{Name: "s", Sysctls: types.Mapping{"a": "b"}}, c2qtypes.Version{Major: 4, Minor: 5}, "sysctls"},
		{"dns <4.7", types.ServiceConfig{Name: "s", DNS: types.StringList{"8.8.8.8"}}, c2qtypes.Version{Major: 4, Minor: 6}, "dns"},
		{"dnsopt <4.7", types.ServiceConfig{Name: "s", DNSOpts: []string{"a"}}, c2qtypes.Version{Major: 4, Minor: 6}, "dns_opt"},
		{"stop_grace_period <5.0", types.ServiceConfig{Name: "s", StopGracePeriod: durationPtr(types.Duration(10_000_000_000))}, c2qtypes.Version{Major: 4, Minor: 9}, "stop_grace_period"},
		{"stop_signal <5.2", types.ServiceConfig{Name: "s", StopSignal: "SIGINT"}, c2qtypes.Version{Major: 5, Minor: 1}, "stop_signal"},
		{"pull_policy <4.6", types.ServiceConfig{Name: "s", PullPolicy: "always"}, c2qtypes.Version{Major: 4, Minor: 5}, "pull_policy"},
		{"label=nested <4.6", types.ServiceConfig{Name: "s", SecurityOpt: []string{"label=nested"}}, c2qtypes.Version{Major: 4, Minor: 5}, "security_opt"},
		{"mask <4.6", types.ServiceConfig{Name: "s", SecurityOpt: []string{"mask=/proc"}}, c2qtypes.Version{Major: 4, Minor: 5}, "security_opt"},
		{"unmask <4.6", types.ServiceConfig{Name: "s", SecurityOpt: []string{"unmask=/proc"}}, c2qtypes.Version{Major: 4, Minor: 5}, "security_opt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := c2qtypes.DefaultConfig()
			cfg.PodmanVersion = tt.version
			dirs := Container(tt.svc, cfg)
			if len(dirs) > 0 {
				t.Fatal("expected no directives but got:", dirs)
			}
			if !hasWarning(cfg, "s", tt.field) {
				t.Fatal("expected WarningSkipped for", tt.field)
			}
		})
	}
}

// T1 tests

func TestContainer_Command(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Command: types.ShellCommand{"npm", "start", "--port", "3000"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Exec", "npm start --port 3000")
}

func TestContainer_Environment(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Environment: types.MappingWithEquals{"NODE_ENV": strPtr("production"), "DEBUG": strPtr("app:*")}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Environment", "NODE_ENV=production")
	assertDirective(t, dirs, "Environment", "DEBUG=app:*")
}

func TestContainer_Environment_BareKey(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Environment: types.MappingWithEquals{"BAZ": nil}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Environment", "BAZ")
}

func TestContainer_Environment_BareKeyUnavailable(t *testing.T) {
	svc := types.ServiceConfig{Environment: types.MappingWithEquals{"BAZ": nil}}
	cfg := c2qtypes.DefaultConfig()
	cfg.PodmanVersion = c2qtypes.Version{Major: 5, Minor: 5}
	dirs := Container(svc, cfg)
	if hasDirective(dirs, "Environment", "BAZ") {
		t.Fatal("expected no bare key Environment directive at 5.5")
	}
	if !hasWarning(cfg, "", "environment") {
		t.Fatal("expected warning for bare key environment")
	}
}

func TestContainer_EnvFile(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", EnvFiles: []types.EnvFile{{Path: "/etc/env", Required: true}, {Path: ".env.local", Required: false}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "EnvironmentFile", "/etc/env")
	assertDirective(t, dirs, "EnvironmentFile", ".env.local")
}

func TestContainer_Tmpfs(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Tmpfs: types.StringList{"/run"}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Tmpfs", "/run")
}

func TestContainer_Logging(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Logging: &types.LoggingConfig{Driver: "json-file", Options: types.Options{"max-size": "10m"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "LogDriver", "json-file")
	assertDirective(t, dirs, "LogOpt", "max-size=10m")
}

func TestContainer_LoggingOptionsUnavailable(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Logging: &types.LoggingConfig{Options: types.Options{"max-size": "10m"}}}
	cfg := c2qtypes.DefaultConfig()
	cfg.PodmanVersion = c2qtypes.Version{Major: 5, Minor: 1}
	dirs := Container(svc, cfg)
	if hasDirective(dirs, "LogOpt", "max-size=10m") {
		t.Fatal("LogOpt should not be emitted before 5.2")
	}
	if !hasWarning(cfg, "web", "logging.options") {
		t.Fatal("expected warning for LogOpt at <5.2")
	}
}

func TestContainer_LogDriverFallback(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", LogDriver: "syslog"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "LogDriver", "syslog")
}

func TestContainer_PidsLimit(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", PidsLimit: 100}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "PidsLimit", "100")
}

func TestContainer_Ulimits(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Ulimits: map[string]*types.UlimitsConfig{"nofile": {Soft: 1024, Hard: 2048}, "nproc": {Single: 65535}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Ulimit", "nofile=1024:2048")
	assertDirective(t, dirs, "Ulimit", "nproc=65535")
}

func TestContainer_ExtraHosts(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", ExtraHosts: types.HostsList{"host.internal": {"10.0.0.1"}, "db": {"192.168.1.5", "192.168.1.6"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "AddHost", "host.internal:10.0.0.1")
	assertDirective(t, dirs, "AddHost", "db:192.168.1.5")
	assertDirective(t, dirs, "AddHost", "db:192.168.1.6")
}

func TestContainer_Volumes(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Volumes: []types.ServiceVolumeConfig{{Type: "bind", Source: "/host/data", Target: "/container/data"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Volume", "/host/data:/container/data:rw")
}

func TestContainer_Volumes_ReadOnly(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Volumes: []types.ServiceVolumeConfig{{Type: "bind", Source: "/host/data", Target: "/app", ReadOnly: true}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Volume", "/host/data:/app:ro")
}

func TestContainer_Ports_ShortSyntax(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Ports: []types.ServicePortConfig{{Published: "8080", Target: 80}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "PublishPort", "8080:80")
}

func TestContainer_Ports_WithProtocol(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Ports: []types.ServicePortConfig{{Published: "53", Target: 53, Protocol: "udp"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "PublishPort", "53:53/udp")
}

func TestContainer_Ports_WithHostIP(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Ports: []types.ServicePortConfig{{HostIP: "127.0.0.1", Published: "8080", Target: 80}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "PublishPort", "127.0.0.1:8080:80")
}

func TestContainer_Ports_TCPOmitted(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Ports: []types.ServicePortConfig{{Published: "8080", Target: 80, Protocol: "tcp"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "PublishPort", "8080:80")
}

func TestContainer_NetworkMode_Host(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", NetworkMode: "host"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Network", "host")
}

func TestContainer_NetworkMode_None(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", NetworkMode: "none"}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Network", "none")
}

func TestContainer_Networks(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Networks: map[string]*types.ServiceNetworkConfig{"frontend": {}, "backend": {}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "Network", "frontend.network")
	assertDirective(t, dirs, "Network", "backend.network")
}

func TestContainer_NetworkAliases(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", Networks: map[string]*types.ServiceNetworkConfig{"frontend": {Aliases: []string{"app", "www"}}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Container(svc, cfg)
	assertDirective(t, dirs, "NetworkAlias", "app:frontend")
	assertDirective(t, dirs, "NetworkAlias", "www:frontend")
}

// Unit tests

func TestUnit_DependsOn(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DependsOn: types.DependsOnConfig{"db": {Condition: "service_started", Required: true}, "redis": {Condition: "service_started", Required: false}}}
	dirs := Unit(svc)

	assertDirective(t, dirs, "Requires", "db.container")
	assertDirective(t, dirs, "Wants", "redis.container")
	assertDirective(t, dirs, "After", "db.container")
	assertDirective(t, dirs, "After", "redis.container")
}

func TestUnit_DependsOn_Restart(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DependsOn: types.DependsOnConfig{"db": {Condition: "service_started", Required: true, Restart: true}}}
	dirs := Unit(svc)

	assertDirective(t, dirs, "Requires", "db.container")
	assertDirective(t, dirs, "BindsTo", "db.container")
	assertDirective(t, dirs, "After", "db.container")
}

func TestUnit_DependsOn_Empty(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", DependsOn: types.DependsOnConfig{}}
	dirs := Unit(svc)
	if len(dirs) != 0 {
		t.Fatalf("expected no directives, got %v", dirs)
	}
}

// Healthcheck tests

func TestHealthcheck_CMD(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: &types.HealthCheckConfig{Test: types.HealthCheckTest{"CMD", "curl", "localhost"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	assertDirective(t, dirs, "HealthCmd", "curl localhost")
}

func TestHealthcheck_CMDSHELL(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: &types.HealthCheckConfig{Test: types.HealthCheckTest{"CMD-SHELL", "curl -f http://localhost || exit 1"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	assertDirective(t, dirs, "HealthCmd", "/bin/sh -c curl -f http://localhost || exit 1")
}

func TestHealthcheck_NONE(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: &types.HealthCheckConfig{Test: types.HealthCheckTest{"NONE"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	if len(dirs) > 0 {
		t.Fatalf("expected no directives for NONE test, got %v", dirs)
	}
}

func TestHealthcheck_Disable(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: &types.HealthCheckConfig{Disable: true, Test: types.HealthCheckTest{"CMD", "curl", "localhost"}}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	if len(dirs) > 0 {
		t.Fatalf("expected no directives when disabled, got %v", dirs)
	}
}

func TestHealthcheck_Interval(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: &types.HealthCheckConfig{Test: types.HealthCheckTest{"CMD", "true"}, Interval: durationPtr(types.Duration(10_000_000_000)), Timeout: durationPtr(types.Duration(5_000_000_000)), Retries: uintPtr(3)}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	assertDirective(t, dirs, "HealthInterval", "10s")
	assertDirective(t, dirs, "HealthTimeout", "5s")
	assertDirective(t, dirs, "HealthRetries", "3")
}

func TestHealthcheck_StartPeriod(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: &types.HealthCheckConfig{Test: types.HealthCheckTest{"CMD", "true"}, StartPeriod: durationPtr(types.Duration(30_000_000_000)), StartInterval: durationPtr(types.Duration(2_000_000_000))}}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	assertDirective(t, dirs, "HealthStartPeriod", "30s")
	assertDirective(t, dirs, "HealthStartupInterval", "2s")
}

func TestHealthcheck_Nil(t *testing.T) {
	svc := types.ServiceConfig{Name: "web", HealthCheck: nil}
	cfg := c2qtypes.DefaultConfig()
	dirs := Healthcheck(svc, cfg)
	if len(dirs) > 0 {
		t.Fatalf("expected no directives for nil healthcheck, got %v", dirs)
	}
}

// Helpers

func assertDirective(t *testing.T, dirs []c2qtypes.Directive, key, value string) {
	t.Helper()
	for _, d := range dirs {
		if d.Key != key {
			continue
		}
		if len(d.Values) == 0 && value == "" {
			return
		}
		for _, v := range d.Values {
			if v == value {
				return
			}
		}
	}
	t.Fatalf("directive %s=%s not found in %v", key, value, dirs)
}

func hasDirective(dirs []c2qtypes.Directive, key, value string) bool {
	for _, d := range dirs {
		if d.Key != key {
			continue
		}
		if len(d.Values) == 0 && value == "" {
			return true
		}
		for _, v := range d.Values {
			if v == value {
				return true
			}
		}
	}
	return false
}

func hasWarning(cfg *c2qtypes.Config, service, field string) bool {
	for _, w := range cfg.Warnings {
		if w.Service == service && w.Field == field {
			return true
		}
	}
	return false
}

func boolPtr(b bool) *bool             { return &b }
func strPtr(s string) *string          { return &s }
func uintPtr(u uint64) *uint64         { return &u }
func durationPtr(d types.Duration) *types.Duration { return &d }
