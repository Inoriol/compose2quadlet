package mapper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	c2qtypes "github.com/inoriol/compose2quadlet/internal/types"
)

func Container(svc types.ServiceConfig, cfg *c2qtypes.Config) []c2qtypes.Directive {
	var dirs []c2qtypes.Directive
	dirs = append(dirs, t0Container(svc, cfg)...)
	dirs = append(dirs, t1Container(svc, cfg)...)
	return dirs
}

func t0Container(svc types.ServiceConfig, cfg *c2qtypes.Config) []c2qtypes.Directive {
	var dirs []c2qtypes.Directive

	if svc.Image != "" {
		dirs = append(dirs, c2qtypes.Directive{Key: "Image", Values: []string{svc.Image}})
	}
	if svc.ContainerName != "" {
		dirs = append(dirs, c2qtypes.Directive{Key: "ContainerName", Values: []string{svc.ContainerName}})
	}
	if svc.Name != "" {
		if cfg.PodmanVersion.AtLeast(5, 3) {
			dirs = append(dirs, c2qtypes.Directive{Key: "ServiceName", Values: []string{svc.Name}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "service_name",
				Message: "requires podman >= 5.3.0",
				Since:   "5.3.0",
			})
		}
	}
	if len(svc.Entrypoint) > 0 {
		if cfg.PodmanVersion.AtLeast(5, 0) {
			dirs = append(dirs, c2qtypes.Directive{Key: "Entrypoint", Values: []string{strings.Join(svc.Entrypoint, " ")}})
		} else {
			dirs = append(dirs, c2qtypes.Directive{Key: "PodmanArgs", Values: []string{"--entrypoint " + strings.Join(svc.Entrypoint, " ")}})
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningDegraded,
				Service: svc.Name,
				Field:   "entrypoint",
				Message: "using PodmanArgs fallback; upgrade to podman >= 5.0.0 for native Entrypoint= support",
				Since:   "5.0.0",
			})
		}
	}
	if svc.Init != nil && *svc.Init {
		dirs = append(dirs, c2qtypes.Directive{Key: "RunInit", Values: []string{"true"}})
	}
	if svc.ReadOnly {
		dirs = append(dirs, c2qtypes.Directive{Key: "ReadOnly", Values: []string{"true"}})
	}
	if svc.User != "" {
		dirs = append(dirs, c2qtypes.Directive{Key: "User", Values: []string{svc.User}})
	}
	if svc.GroupAdd != nil {
		if cfg.PodmanVersion.AtLeast(5, 1) {
			for _, g := range svc.GroupAdd {
				dirs = append(dirs, c2qtypes.Directive{Key: "GroupAdd", Values: []string{g}})
			}
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "group_add",
				Message: "requires podman >= 5.1.0",
				Since:   "5.1.0",
			})
		}
	}
	if svc.WorkingDir != "" {
		if cfg.PodmanVersion.AtLeast(4, 6) {
			dirs = append(dirs, c2qtypes.Directive{Key: "WorkingDir", Values: []string{svc.WorkingDir}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "working_dir",
				Message: "requires podman >= 4.6.0",
				Since:   "4.6.0",
			})
		}
	}

	dirs = append(dirs, SecurityOpts(svc.SecurityOpt, svc.Name, cfg)...)

	if svc.MemLimit > 0 {
		if cfg.PodmanVersion.AtLeast(5, 5) {
			dirs = append(dirs, c2qtypes.Directive{Key: "Memory", Values: []string{strconv.FormatInt(int64(svc.MemLimit), 10)}})
		}
	}
	if svc.Cgroup == "host" {
		if cfg.PodmanVersion.AtLeast(5, 3) {
			dirs = append(dirs, c2qtypes.Directive{Key: "CgroupsMode", Values: []string{"host"}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "cgroup",
				Message: "requires podman >= 5.3.0",
				Since:   "5.3.0",
			})
		}
	}
	if svc.Cgroup == "private" {
		dirs = append(dirs, c2qtypes.Directive{Key: "PodmanArgs", Values: []string{"--cgroupns private"}})
	}

	if svc.UserNSMode != "" {
		if cfg.PodmanVersion.AtLeast(4, 5) {
			dirs = append(dirs, c2qtypes.Directive{Key: "UserNS", Values: []string{svc.UserNSMode}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "userns_mode",
				Message: "requires podman >= 4.5.0",
				Since:   "4.5.0",
			})
		}
	}
	if svc.Hostname != "" {
		if cfg.PodmanVersion.AtLeast(4, 6) {
			dirs = append(dirs, c2qtypes.Directive{Key: "HostName", Values: []string{svc.Hostname}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "hostname",
				Message: "requires podman >= 4.6.0",
				Since:   "4.6.0",
			})
		}
	}
	if svc.ShmSize != 0 {
		if cfg.PodmanVersion.AtLeast(4, 7) {
			dirs = append(dirs, c2qtypes.Directive{Key: "ShmSize", Values: []string{strconv.FormatInt(int64(svc.ShmSize), 10)}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "shm_size",
				Message: "requires podman >= 4.7.0",
				Since:   "4.7.0",
			})
		}
	}
	if len(svc.Sysctls) > 0 {
		if cfg.PodmanVersion.AtLeast(4, 6) {
			for k, v := range svc.Sysctls {
				dirs = append(dirs, c2qtypes.Directive{Key: "Sysctl", Values: []string{k + "=" + v}})
			}
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "sysctls",
				Message: "requires podman >= 4.6.0",
				Since:   "4.6.0",
			})
		}
	}
	if len(svc.DNS) > 0 {
		if cfg.PodmanVersion.AtLeast(4, 7) {
			for _, d := range svc.DNS {
				dirs = append(dirs, c2qtypes.Directive{Key: "DNS", Values: []string{d}})
			}
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "dns",
				Message: "requires podman >= 4.7.0",
				Since:   "4.7.0",
			})
		}
	}
	if len(svc.DNSSearch) > 0 {
		if cfg.PodmanVersion.AtLeast(4, 7) {
			for _, d := range svc.DNSSearch {
				dirs = append(dirs, c2qtypes.Directive{Key: "DNSSearch", Values: []string{d}})
			}
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "dns_search",
				Message: "requires podman >= 4.7.0",
				Since:   "4.7.0",
			})
		}
	}
	for _, d := range svc.DNSOpts {
		if cfg.PodmanVersion.AtLeast(4, 7) {
			dirs = append(dirs, c2qtypes.Directive{Key: "DNSOption", Values: []string{d}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "dns_opt",
				Message: "requires podman >= 4.7.0",
				Since:   "4.7.0",
			})
			break
		}
	}
	if svc.StopGracePeriod != nil {
		if cfg.PodmanVersion.AtLeast(5, 0) {
			dirs = append(dirs, c2qtypes.Directive{Key: "StopTimeout", Values: []string{time.Duration(*svc.StopGracePeriod).String()}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "stop_grace_period",
				Message: "requires podman >= 5.0.0",
				Since:   "5.0.0",
			})
		}
	}
	if svc.StopSignal != "" {
		if cfg.PodmanVersion.AtLeast(5, 2) {
			dirs = append(dirs, c2qtypes.Directive{Key: "StopSignal", Values: []string{svc.StopSignal}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "stop_signal",
				Message: "requires podman >= 5.2.0",
				Since:   "5.2.0",
			})
		}
	}
	if svc.PullPolicy != "" {
		if cfg.PodmanVersion.AtLeast(4, 6) {
			dirs = append(dirs, c2qtypes.Directive{Key: "Pull", Values: []string{svc.PullPolicy}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "pull_policy",
				Message: "requires podman >= 4.6.0",
				Since:   "4.6.0",
			})
		}
	}

	for k, v := range svc.Labels {
		dirs = append(dirs, c2qtypes.Directive{Key: "Label", Values: []string{fmt.Sprintf("%s=%s", k, v)}})
	}
	for k, v := range svc.Annotations {
		dirs = append(dirs, c2qtypes.Directive{Key: "Annotation", Values: []string{fmt.Sprintf("%s=%s", k, v)}})
	}
	for _, c := range svc.CapAdd {
		dirs = append(dirs, c2qtypes.Directive{Key: "AddCapability", Values: []string{c}})
	}
	for _, c := range svc.CapDrop {
		dirs = append(dirs, c2qtypes.Directive{Key: "DropCapability", Values: []string{c}})
	}
	for _, e := range svc.Expose {
		dirs = append(dirs, c2qtypes.Directive{Key: "ExposeHostPort", Values: []string{e}})
	}

	return dirs
}

func t1Container(svc types.ServiceConfig, cfg *c2qtypes.Config) []c2qtypes.Directive {
	var dirs []c2qtypes.Directive

	if len(svc.Command) > 0 {
		dirs = append(dirs, c2qtypes.Directive{Key: "Exec", Values: []string{strings.Join(svc.Command, " ")}})
	}

	for k, v := range svc.Environment {
		if v != nil {
			dirs = append(dirs, c2qtypes.Directive{Key: "Environment", Values: []string{k + "=" + *v}})
		} else {
			if cfg.PodmanVersion.AtLeast(5, 6) {
				dirs = append(dirs, c2qtypes.Directive{Key: "Environment", Values: []string{k}})
			} else {
				cfg.Warn(c2qtypes.Warning{
					Level:   c2qtypes.WarningSkipped,
					Service: svc.Name,
					Field:   "environment",
					Message: "bare key requires podman >= 5.6.0",
					Since:   "5.6.0",
				})
			}
		}
	}
	for _, f := range svc.EnvFiles {
		dirs = append(dirs, c2qtypes.Directive{Key: "EnvironmentFile", Values: []string{f.Path}})
	}
	for _, t := range svc.Tmpfs {
		dirs = append(dirs, c2qtypes.Directive{Key: "Tmpfs", Values: []string{t}})
	}

	if svc.Logging != nil {
		if svc.Logging.Driver != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "LogDriver", Values: []string{svc.Logging.Driver}})
		}
		if len(svc.Logging.Options) > 0 {
			if cfg.PodmanVersion.AtLeast(5, 2) {
				for k, v := range svc.Logging.Options {
					dirs = append(dirs, c2qtypes.Directive{Key: "LogOpt", Values: []string{k + "=" + v}})
				}
			} else {
				cfg.Warn(c2qtypes.Warning{
					Level:   c2qtypes.WarningSkipped,
					Service: svc.Name,
					Field:   "logging.options",
					Message: "requires podman >= 5.2.0",
					Since:   "5.2.0",
				})
			}
		}
	} else {
		if svc.LogDriver != "" {
			dirs = append(dirs, c2qtypes.Directive{Key: "LogDriver", Values: []string{svc.LogDriver}})
		}
		for k, v := range svc.LogOpt {
			if cfg.PodmanVersion.AtLeast(5, 2) {
				dirs = append(dirs, c2qtypes.Directive{Key: "LogOpt", Values: []string{k + "=" + v}})
			} else {
				cfg.Warn(c2qtypes.Warning{
					Level:   c2qtypes.WarningSkipped,
					Service: svc.Name,
					Field:   "logging.options",
					Message: "requires podman >= 5.2.0",
					Since:   "5.2.0",
				})
				break
			}
		}
	}

	if svc.PidsLimit != 0 {
		if cfg.PodmanVersion.AtLeast(4, 7) {
			dirs = append(dirs, c2qtypes.Directive{Key: "PidsLimit", Values: []string{strconv.FormatInt(svc.PidsLimit, 10)}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "pids_limit",
				Message: "requires podman >= 4.7.0",
				Since:   "4.7.0",
			})
		}
	}

	for name, u := range svc.Ulimits {
		if !cfg.PodmanVersion.AtLeast(4, 7) {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "ulimits",
				Message: "requires podman >= 4.7.0",
				Since:   "4.7.0",
			})
			break
		}
		var val string
		if u.Single > 0 {
			val = strconv.Itoa(u.Single)
		} else {
			val = strconv.Itoa(u.Soft) + ":" + strconv.Itoa(u.Hard)
		}
		dirs = append(dirs, c2qtypes.Directive{Key: "Ulimit", Values: []string{name + "=" + val}})
	}

	if len(svc.ExtraHosts) > 0 {
		if cfg.PodmanVersion.AtLeast(5, 3) {
			for host, ips := range svc.ExtraHosts {
				for _, ip := range ips {
					dirs = append(dirs, c2qtypes.Directive{Key: "AddHost", Values: []string{host + ":" + ip}})
				}
			}
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "extra_hosts",
				Message: "requires podman >= 5.3.0",
				Since:   "5.3.0",
			})
		}
	}

	for _, v := range svc.Volumes {
		if v.Type == types.VolumeTypeBind {
			dirs = append(dirs, c2qtypes.Directive{Key: "Mount", Values: []string{formatBindMount(v)}})
		} else if v.Type == types.VolumeTypeTmpfs {
			dirs = append(dirs, c2qtypes.Directive{Key: "Tmpfs", Values: []string{formatTmpfsMount(v)}})
		} else {
			dirs = append(dirs, c2qtypes.Directive{Key: "Volume", Values: []string{v.String()}})
		}
	}

	for _, p := range svc.Ports {
		dirs = append(dirs, c2qtypes.Directive{Key: "PublishPort", Values: []string{formatPort(p)}})
	}

	if svc.NetworkMode == "host" {
		dirs = append(dirs, c2qtypes.Directive{Key: "Network", Values: []string{"host"}})
	} else if svc.NetworkMode == "none" {
		dirs = append(dirs, c2qtypes.Directive{Key: "Network", Values: []string{"none"}})
	} else if strings.HasPrefix(svc.NetworkMode, "service:") {
		target := strings.TrimPrefix(svc.NetworkMode, "service:")
		if cfg.PodmanVersion.AtLeast(5, 3) {
			dirs = append(dirs, c2qtypes.Directive{Key: "Network", Values: []string{"container:" + target + ".container"}})
		} else {
			cfg.Warn(c2qtypes.Warning{
				Level:   c2qtypes.WarningSkipped,
				Service: svc.Name,
				Field:   "network_mode",
				Message: "network_mode: service:<name> requires podman >= 5.3.0",
				Since:   "5.3.0",
			})
		}
	} else if svc.NetworkMode == "" || svc.NetworkMode == "bridge" {
		for name, net := range svc.Networks {
			dirs = append(dirs, c2qtypes.Directive{Key: "Network", Values: []string{name + ".network"}})
			if net.Ipv4Address != "" {
				dirs = append(dirs, c2qtypes.Directive{Key: "IP", Values: []string{net.Ipv4Address}})
			}
			if net.Ipv6Address != "" {
				dirs = append(dirs, c2qtypes.Directive{Key: "IP6", Values: []string{net.Ipv6Address}})
			}
			if len(net.Aliases) > 0 {
				if cfg.PodmanVersion.AtLeast(5, 2) {
					for _, alias := range net.Aliases {
						dirs = append(dirs, c2qtypes.Directive{Key: "NetworkAlias", Values: []string{alias + ":" + name}})
					}
				} else {
					cfg.Warn(c2qtypes.Warning{
						Level:   c2qtypes.WarningSkipped,
						Service: svc.Name,
						Field:   "networks." + name + ".aliases",
						Message: "requires podman >= 5.2.0",
						Since:   "5.2.0",
					})
				}
			}
		}
	}

	return dirs
}

func formatBindMount(v types.ServiceVolumeConfig) string {
	var parts []string
	parts = append(parts, "type=bind")
	parts = append(parts, "source="+v.Source)
	parts = append(parts, "destination="+v.Target)
	if v.ReadOnly {
		parts = append(parts, "readonly")
	}
	if v.Bind != nil {
		if v.Bind.Propagation != "" {
			parts = append(parts, "bind-propagation="+v.Bind.Propagation)
		}
		if v.Bind.SELinux != "" {
			parts = append(parts, "selinux="+v.Bind.SELinux)
		}
	}
	return strings.Join(parts, ",")
}

func formatTmpfsMount(v types.ServiceVolumeConfig) string {
	val := v.Target
	if v.Tmpfs != nil {
		var opts []string
		if v.Tmpfs.Size > 0 {
			opts = append(opts, fmt.Sprintf("size=%d", v.Tmpfs.Size))
		}
		if v.Tmpfs.Mode > 0 {
			opts = append(opts, fmt.Sprintf("mode=%o", v.Tmpfs.Mode))
		}
		if len(opts) > 0 {
			val += ":" + strings.Join(opts, ",")
		}
	}
	return val
}
