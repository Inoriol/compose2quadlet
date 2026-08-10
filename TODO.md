# Implementation TODO — ordered easiest to hardest

Scale: **T0 (trivial)** → **T1 (easy)** → **T2 (medium)** → **T3 (hard)** → 
**T4 (veryhard)** → **T5 (structural)**. Within each tier, fields are grouped by 
the mapper file where they will be implemented.

---

## T0 — Trivial (1:1 string passthrough, baseline 4.4.0)

Zero transformation. Compose field value copied verbatim to a quadlet directive.

Status: ✅ = done, ❌ = no compose-go v2 field (commentary inline), ⏳ = moved to higher tier

File: `mapper/container.go`

| # | Compose field | → Directive | Since | Status |
|---|---|---|---|---|
| 1 | `image` | `Image=` | 4.4.0 | ✅ |
| 2 | `init` | `RunInit=` | 4.4.0 | ✅ |
| 3 | `read_only` | `ReadOnly=` | 4.4.0 | ✅ |
| 4 | `container_name` | `ContainerName=` | 4.4.0 | ✅ |
| 5 | `labels` (map/list) | `Label=` | 4.4.0 | ✅ |
| 6 | `annotations` (map/list) | `Annotation=` | 4.4.0 | ✅ |
| 7 | `cap_add` | `AddCapability=` | 4.4.0 | ✅ |
| 8 | `cap_drop` | `DropCapability=` | 4.4.0 | ✅ |
| 9 | `group` | `Group=` | 4.4.0 | ❌ compose-go uses `group_add`, not `group` |
| 10 | `user` | `User=` | 4.4.0 | ✅ |
| 11 | `expose` | `ExposeHostPort=` | 4.4.0 | ✅ |

File: `mapper/security.go`

| # | Compose field | → Directive | Since | Status |
|---|---|---|---|---|
| 12 | `security_opt: no-new-privileges` | `NoNewPrivileges=` | 4.4.0 | ✅ |
| 13 | `security_opt: label=disable` | `SecurityLabelDisable=` | 4.4.0 | ✅ |
| 14 | `security_opt: label=nested` | `SecurityLabelNested=` | 4.6.0 | ✅ |
| 15 | `security_opt: seccomp=<path>` | `SeccompProfile=` | 4.4.0 | ✅ |
| 16 | `security_opt: label=type:<t>` | `SecurityLabelType=` | 4.4.0 | ✅ |
| 17 | `security_opt: label=level:<l>` | `SecurityLabelLevel=` | 4.4.0 | ✅ |
| 18 | `security_opt: label=filetype:<ft>` | `SecurityLabelFileType=` | 4.4.0 | ✅ |
| 19 | `security_opt: mask=<path>` | `Mask=` | 4.6.0 | ✅ |
| 20 | `security_opt: unmask=<path>` | `Unmask=` | 4.6.0 | ✅ |
| 21 | `security_opt: apparmor=<p>` | `AppArmor=` | 5.8.0 | ✅ |
| 22 | `userns_mode` | `UserNS=` | 4.5.0 | ✅ |
| 23 | `group_add` | `GroupAdd=` | 5.1.0 | ✅ |

File: `mapper/container.go`

| # | Compose field | → Directive | Since | Status |
|---|---|---|---|---|
| 24 | `uid_map` | `UIDMap=` | 4.8.0 | ❌ no field in compose-go v2 — use PodmanArgs= |
| 25 | `gid_map` | `GIDMap=` | 4.8.0 | ❌ no field in compose-go v2 — use PodmanArgs= |
| 26 | `sub_uid_map` | `SubUIDMap=` | 4.8.0 | ❌ no field in compose-go v2 — use PodmanArgs= |
| 27 | `sub_gid_map` | `SubGIDMap=` | 4.8.0 | ❌ no field in compose-go v2 — use PodmanArgs= |
| 28 | `read_only` (+ tmpfs flag) | `ReadOnlyTmpfs=` | 4.8.0 | ❌ ReadOnly is bool, no separate tmpfs distinction in compose-go |
| 29 | `devices` (`HOST:CONTAINER[:PERMS]`) | `AddDevice=` | 4.4.0 | ⏳ exists as `Devices []DeviceMapping` — needs struct parsing (T1/T2) |
| 30 | `devices` (CDI syntax) | `AddDevice=` | 4.4.0 | ⏳ see #29 |
| 31 | `working_dir` | `WorkingDir=` | 4.6.0 | ✅ |
| 32 | `hostname` | `HostName=` | 4.6.0 | ✅ |
| 33 | `shm_size` | `ShmSize=` | 4.7.0 | ✅ |
| 34 | `sysctls` (map/list) | `Sysctl=` | 4.6.0 | ✅ |
| 35 | `dns` | `DNS=` | 4.7.0 | ✅ |
| 36 | `dns_search` | `DNSSearch=` | 4.7.0 | ✅ |
| 37 | `dns_opt` | `DNSOption=` | 4.7.0 | ✅ |
| 38 | `stop_grace_period` | `StopTimeout=` | 5.0.0 | ✅ |
| 39 | `stop_signal` | `StopSignal=` | 5.2.0 | ✅ |
| 40 | `pull_policy` | `Pull=` | 4.6.0 | ✅ |
| 41 | `notify` | `Notify=` | 5.0.0 | ❌ no field in compose-go v2 — derive from healthcheck/deploy? |
| 42 | `timezone` | `Timezone=` | 4.6.0 | ❌ no field in compose-go v2 — may set via environment TZ= |
| 43 | `rootfs` | `Rootfs=` | 4.5.0 | ❌ no field in compose-go v2 |
| 44 | `reload_signal` | `ReloadSignal=` | 5.5.0 | ❌ no field in compose-go v2 |
| 45 | `image_volume` | `ImageVolume=` | 6.1.0 | ❌ no field in compose-go v2 |
| 46 | `service_name` | `ServiceName=` | 5.3.0 | ✅ mapped from `svc.Name` |
| 47 | `containers_conf_module` | `ContainersConfModule=` | ? | ❌ no field in compose-go v2 |
| 48 | `http_proxy` | `HttpProxy=` | 5.7.0 | ❌ no field in compose-go v2 |
| 49 | `environment_host` | `EnvironmentHost=` | ? | ❌ no field in compose-go v2 |

File: `mapper/container.go` — network IP (moved from network.go)

| # | Compose field | → Directive | Since | Status |
|---|---|---|---|---|
| 50 | `networks.<name>.ipv4_address` | `IP=` | 4.5.0 | ✅ |
| 51 | `networks.<name>.ipv6_address` | `IP6=` | 4.5.0 | ✅ |

---

## T1 — Easy (minor formatting / type conversion)

Simple transforms: list→multi-value, map→key=value, append suffix, path handling.

Status: ✅ = done, ❌ = no compose-go v2 field

File: `mapper/container.go`

| # | Compose field | → Directive | Notes | Since | Status |
|---|---|---|---|---|---|
| 52 | `command` (string/list) | `Exec=` | list→space-joined string | 4.4.0 | ✅ |
| 53 | `environment` (map) | `Environment=` | map→`KEY=VALUE` per line | 4.4.0 | ✅ |
| 54 | `environment` (list) | `Environment=` | `KEY=VALUE` already | 4.4.0 | ✅ same as #53 |
| 55 | `env_file` (string/list) | `EnvironmentFile=` | multi-value | 4.4.0 | ✅ |
| 56 | `env_file` (`required: false`) | `EnvironmentFile=` | same, ignore required | 4.4.0 | ✅ |
| 57 | `tmpfs` (string/long) | `Tmpfs=` | format: `path:opts` | 4.5.0 | ✅ |
| 58 | `logging.driver` | `LogDriver=` | passthrough | 4.5.0 | ✅ |
| 59 | `logging.options` | `LogOpt=` | map→`key=value` | 5.2.0 | ✅ |
| 60 | `pids_limit` | `PidsLimit=` | passthrough | 4.7.0 | ✅ |
| 61 | `ulimits` | `Ulimit=` | format: `name=soft:hard` | 4.7.0 | ✅ |
| 62 | `extra_hosts` (list/map) | `AddHost=` | `host:ip` format | 5.3.0 | ✅ |
| 63 | `reload_cmd` | `ReloadCmd=` | list→space-joined | 5.5.0 | ❌ no field in compose-go v2 |
| 64 | `environment` (key-only) | `Environment=` | emit bare key | 5.6.0 | ✅ |

File: `mapper/container.go` — volume mounts

| # | Compose field | → Directive | Notes | Since | Status |
|---|---|---|---|---|---|
| 65 | `volumes` (short syntax) | `Volume=` | `src:dst:opts` passthrough | 4.4.0 | ✅ via ServiceVolumeConfig.String() |
| 66 | `volumes.read_only` | `Volume=... :ro` | append `:ro` | 4.4.0 | ✅ via String() |
| 67 | `volumes.selinux` | `Volume=... :z`/`:Z` | append selinux label | 4.4.0 | ⏳ compose compiles but String() preserves inline flag |
| 68 | `volumes.nocopy` | `Volume=... :nocopy` | append `:nocopy` | 4.4.0 | ✅ via String() |

File: `mapper/container.go` — ports

| # | Compose field | → Directive | Notes | Since | Status |
|---|---|---|---|---|---|
| 69 | `ports` (short syntax) | `PublishPort=` | parse `HOST:CONTAINER[/proto]` | 4.4.0 | ✅ |
| 70 | `ports` (long syntax) | `PublishPort=` | `target/pub/proto/host_ip` | 4.4.0 | ✅ |

File: `mapper/container.go` — network modes

| # | Compose field | → Directive | Notes | Since | Status |
|---|---|---|---|---|---|
| 71 | `network_mode: host` | `Network=host` | literal | 4.4.0 | ✅ |
| 72 | `network_mode: none` | `Network=none` | literal | 4.4.0 | ✅ |
| 73 | `networks` (list) | `Network=` | emit `<name>.network` per entry | 4.4.0 | ✅ |
| 74 | `networks.aliases` | `NetworkAlias=` | multi-value, one per alias | 5.2.0 | ✅ |
| 50 | `networks.<name>.ipv4_address` | `IP=` | per-network IP | 4.5.0 | ✅ |
| 51 | `networks.<name>.ipv6_address` | `IP6=` | per-network IP | 4.5.0 | ✅ |

File: `mapper/container.go` — depends_on

| # | Compose field | → Directive | Notes | Since | Status |
|---|---|---|---|---|---|
| 75 | `depends_on` (short list) | `After=` / `Requires=` | `<name>.container` per dep | 4.4.0 | ✅ |
| 76 | `depends_on` (long, `condition: service_started`) | `After=` / `Requires=` | default condition | 4.4.0 | ✅ |
| 77 | `depends_on` (long, `required: false`) | `Wants=` | instead of `Requires=` | 4.4.0 | ✅ |
| 78 | `depends_on` (long, `restart: true`) | `BindsTo=` | additional binding | 4.4.0 | ✅ |

File: `mapper/container.go` — healthcheck

| # | Compose field | → Directive | Notes | Since | Status |
|---|---|---|---|---|---|
| 79 | `healthcheck.test` (CMD) | `HealthCmd=` | list→space-joined | 4.5.0 | ✅ |
| 80 | `healthcheck.test` (CMD-SHELL) | `HealthCmd=` | wrap `/bin/sh -c` | 4.5.0 | ✅ |
| 81 | `healthcheck.interval` | `HealthInterval=` | passthrough | 4.5.0 | ✅ |
| 82 | `healthcheck.timeout` | `HealthTimeout=` | passthrough | 4.5.0 | ✅ |
| 83 | `healthcheck.retries` | `HealthRetries=` | passthrough | 4.5.0 | ✅ |
| 84 | `healthcheck.start_period` | `HealthStartPeriod=` | passthrough | 4.5.0 | ✅ |
| 85 | `healthcheck.start_interval` | `HealthStartupInterval=` | passthrough | 4.5.0 | ✅ |
| 94 | `healthcheck.test` (NONE) | *(omit health directives)* | skip all health directives | 4.5.0 | ✅ |
| 95 | `healthcheck.disable: true` | *(omit health directives)* | skip all health directives | 4.5.0 | ✅ |
| 86 | `healthcheck.on_failure` | `HealthOnFailure=` | passthrough | 4.5.0 | ❌ not in compose-go HealthCheckConfig |
| 87 | `healthcheck.log_destination` | `HealthLogDestination=` | passthrough | 5.3.0 | ❌ not in compose-go HealthCheckConfig |
| 88 | `healthcheck.max_log_count` | `HealthMaxLogCount=` | passthrough | 5.3.0 | ❌ not in compose-go HealthCheckConfig |
| 89 | `healthcheck.max_log_size` | `HealthMaxLogSize=` | passthrough | 5.3.0 | ❌ not in compose-go HealthCheckConfig |

---

## T2 — Medium (complex formatting, multi-value logic, version gating)

Fields that need structured parsing, emit multiple related directives, or have
version-conditional P1 vs P3 fallback paths.

File: `mapper/container.go`

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 90 | `entrypoint` (string/list) | `Entrypoint=` | list→space-joined; P3 PodmanArgs fallback pre-5.0 | 5.0.0 |
| 91 | `volumes` (long, bind) | `Mount=type=bind,...` | construct mount string | 4.5.0 |
| 92 | `volumes` (long, volume) | `Volume=` | long→short conversion | 4.4.0 |
| 93 | `tmpfs` (long syntax, `size`/`mode`) | `Tmpfs=` | construct `path:opt1,opt2` | 4.5.0 |
| 94 | `healthcheck.test` (NONE) | *(omit health directives)* | skip all health directives | 4.5.0 | ⏳ moved to T1 ✅ |
| 95 | `healthcheck.disable: true` | *(omit health directives)* | skip all health directives | 4.5.0 | ⏳ moved to T1 ✅ |

File: `mapper/container.go` — resource version gate (P1 quadlet vs P2 systemd)

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 96 | `mem_limit` (P1) | `Memory=` | P1 since 5.5.0; P2 fallback `MemoryMax=[Service]` | 5.5.0 |
| 97 | `mem_limit` (P2 alt) | `MemoryMax=` [Service] | pre-5.5.0 preferred path | sd 231 |
| 98 | `pids_limit` (P2 alt) | `TasksMax=` [Service] | alt to direct PidsLimit= | sd 227 |
| 99 | `oom_kill_disable` (P2) | `ManagedOOMMemoryPressure=kill` [Service] | default P2; P3 alt `--oom-kill-disable` | sd 247 |
| 100 | `cgroup_parent` (P2) | `Slice=` [Service] | default P2; P3 alt `--cgroup-parent ...` | sd 208 |

File: `mapper/service.go` — systemd resource-control

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 101 | `mem_reservation` | `MemoryLow=` [Service] | | sd 240 |
| 102 | `memswap_limit` | `MemorySwapMax=` [Service] | | sd 232 |
| 103 | `cpus` | `CPUQuota=` [Service] | fractional→percent×period | sd 213 |
| 104 | `cpu_shares` | `CPUWeight=` [Service] | shares→weight conversion | sd 232 |
| 105 | `cpu_period` | `CPUQuotaPeriodSec=` [Service] | | sd 242 |
| 106 | `cpu_quota` | `CPUQuota=` [Service] | | sd 213 |
| 107 | `cpuset` | `AllowedCPUs=` [Service] | | sd 244 |
| 108 | `oom_score_adj` | `OOMScoreAdjust=` [Service] | | sd 208 |
| 109 | `blkio_config.weight` | `IOWeight=` [Service] | | sd 230 |
| 110 | `blkio_config.weight_device` | `IODeviceWeight=` [Service] | | sd 230 |
| 111 | `blkio_config.device_read_bps` | `IOReadBandwidthMax=` [Service] | | sd 230 |
| 112 | `blkio_config.device_write_bps` | `IOWriteBandwidthMax=` [Service] | | sd 230 |
| 113 | `blkio_config.device_read_iops` | `IOReadIOPSMax=` [Service] | | sd 230 |
| 114 | `blkio_config.device_write_iops` | `IOWriteIOPSMax=` [Service] | | sd 230 |

File: `mapper/service.go` — restart logic

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 115 | `restart: no/always` | `Restart=` [Service] | direct mapping | sd 208 |
| 116 | `restart: on-failure` | `Restart=on-failure` [Service] | | sd 208 |
| 117 | `restart: on-failure:<N>` | `Restart=on-failure` + `StartLimitBurst=N` [Service] | split into two directives | sd 208 |
| 118 | `restart: unless-stopped` | `Restart=always` [Service] | compose→systemd semantic mapping | sd 208 |
| 119 | `deploy.restart_policy.condition` | `Restart=` [Service] | same as `restart` but deploy section | sd 208 |
| 120 | `deploy.restart_policy.delay` | `RestartSec=` [Service] | | sd 208 |
| 121 | `deploy.restart_policy.max_attempts` | `StartLimitBurst=` / `StartLimitIntervalSec=` [Service] | | sd 208 |
| 122 | `deploy.restart_policy.window` | `RuntimeMaxSec=` [Service] | | sd 229 |

File: `mapper/service.go` — deploy.resources

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 123 | `deploy.resources.limits.cpus` | `CPUQuota=` [Service] | | sd 213 |
| 124 | `deploy.resources.limits.memory` | `MemoryMax=` [Service] | | sd 231 |
| 125 | `deploy.resources.limits.pids` | `TasksMax=` [Service] | | sd 227 |
| 126 | `deploy.resources.reservations.cpus` | `CPUWeight=` [Service] | | sd 232 |
| 127 | `deploy.resources.reservations.memory` | `MemoryLow=` [Service] | | sd 240 |
| 128 | `ulimits` (P2 alt) | `Limit*= ` [Service] | emit `LimitXXX=value` per ulimit name | sd 208 |

File: `mapper/network.go` — network_mode service

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 129 | `network_mode: service:<name>` | `Network=container:<name>.container` | cross-ref | 5.3.0 |

File: `mapper/container.go` — cgroup modes

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 130 | `cgroup: host` | `CgroupsMode=host` | | 5.3.0 |
| 131 | `cgroup: private` | `PodmanArgs=--cgroupns private` (P3) | | 4.6.0 |

---

## T3 — Hard (PodmanArgs P3, complex parsing, multi-directive)

Fields with no native quadlet directive, requiring Passthrough formatting to
`PodmanArgs=` or `GlobalArgs=`. Also includes complex compose field parsing
(security_opt, devices with options).

File: `mapper/container.go` — P3 passthroughs

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 132 | `tty` | `PodmanArgs=--tty` | flag, no value | 4.6.0 |
| 133 | `stdin_open` | `PodmanArgs=--attach stdin` | flag | 4.6.0 |
| 134 | `runtime` | `GlobalArgs=--runtime <name>` | | 4.6.0 |
| 135 | `mac_address` | `PodmanArgs=--mac-address ...` | | 4.6.0 |
| 136 | `networks.mac_address` | `PodmanArgs=--mac-address ...` | per-network flag | 4.6.0 |
| 137 | `ipc: shareable` | `PodmanArgs=--ipc shareable` | | 4.6.0 |
| 138 | `pid: host` | `PodmanArgs=--pid host` | | 4.6.0 |
| 139 | `uts: host` | `PodmanArgs=--uts host` | | 4.6.0 |
| 140 | `privileged` | `PodmanArgs=--privileged` | flag | 4.6.0 |
| 141 | `mem_swappiness` | `PodmanArgs=--memory-swappiness ...` | | 4.6.0 |
| 142 | `cpu_rt_runtime` | `PodmanArgs=--cpu-rt-runtime ...` | | 4.6.0 |
| 143 | `cpu_rt_period` | `PodmanArgs=--cpu-rt-period ...` | | 4.6.0 |
| 144 | `device_cgroup_rules` | `PodmanArgs=--device-cgroup-rule ...` | multi-value | 4.6.0 |
| 145 | `storage_opt` | `GlobalArgs=--storage-opt ...` | multi-value | 4.6.0 |
| 146 | `oom_kill_disable` (P3 alt) | `PodmanArgs=--oom-kill-disable` | flag | 4.6.0 |
| 147 | `cgroup_parent` (P3 alt) | `PodmanArgs=--cgroup-parent ...` | | 4.6.0 |
| 148 | `entrypoint` (P3 fallback) | `PodmanArgs=--entrypoint ...` | pre-5.0.0 | 4.6.0 |

File: `mapper/unit.go` — complex conditionals

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 149 | `depends_on` (`condition: service_healthy`) | `ExecStartPre=` [Service] | health polling hook | sd 208 |
| 150 | `depends_on` (`condition: service_completed_successfully`) | `After=` / `Requires=` [Unit] | oneshot completion wait | 4.4.0 |

File: `mapper/container.go` — misc hard

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 151 | `label_file` | `Label=` | read file, inline as Label= | — |

---

## T4 — Very hard (secrets/configs pre-intercept, advanced version gating)

Fields handled by the pre-mapping interceptor before field mapping runs. Requires
multiple code paths depending on secret/config type.

File: `mapper/secrets.go`

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 152 | `secrets` (short syntax) | *(pre-mapping)* | route to correct handler | — |
| 153 | `secrets` (long, external) | `Secret=` | podman secret | 4.5.0 |
| 154 | `secrets` (long, file) | `Volume=<path>:/run/secrets/<name>:ro` | bind mount | 4.4.0 |
| 155 | `secrets` (long, environment) | `Volume=<path>:/run/secrets/<name>:ro` | same as file | 4.4.0 |
| 156 | `configs` (short syntax) | `Mount=type=bind,source=<path>,target=/<name>` | bind mount | 4.5.0 |
| 157 | `configs` (long syntax) | `Mount=type=bind,...` | full mount options | 4.5.0 |

File: `mapper/container.go` — version-gated entrypoint path

| # | Compose field | → Directive | Notes | Since |
|---|---|---|---|---|
| 158 | `entrypoint` (version switch) | P1 `Entrypoint=` / P3 `PodmanArgs=--entrypoint` | central version-gate logic | 5.0.0 |

---

## T5 — Structural (generates separate quadlet units)

These fields don't produce in-container directives — they spawn entirely new
`QuadletUnit` structs. The hardest tier because each requires its own mapper file
and cross-unit reference wiring.

File: `mapper/image.go` — `.image` companion unit

| # | Compose field | → Quadlet | Notes | Since |
|---|---|---|---|---|
| 159 | `image` (unit gen) | `.image` unit | one per service with `image:` | 4.8.0 |
| 160 | `image` (Image directive) | `Image=` [Image] | in `.image` unit | 4.8.0 |
| 161 | `pull_policy` (image unit) | `Policy=` [Image] | | 5.6.0 |
| 162 | `platform` OS | `OS=` [Image] | parse `os/arch/variant` | 4.8.0 |
| 163 | `platform` arch | `Arch=` [Image] | | 4.8.0 |
| 164 | `platform` variant | `Variant=` [Image] | | 4.8.0 |
| 165 | retries | `Retry=` [Image] | from config or default | 5.5.0 |
| 166 | retry delay | `RetryDelay=` [Image] | from config or default | 5.5.0 |

File: `mapper/build.go` — `.build` unit

| # | Compose field | → Quadlet | Notes | Since |
|---|---|---|---|---|
| 167 | `build` (unit gen) | `.build` unit | fatal error pre-5.2.0 | 5.2.0 |
| 168 | `build.context` | `SetWorkingDirectory=` [Build] | | 5.2.0 |
| 169 | `build.dockerfile` | `File=` [Build] | | 5.2.0 |
| 170 | `build.args` | `BuildArg=` [Build] | | 5.7.0 |
| 171 | `build.target` | `Target=` [Build] | | 5.2.0 |
| 172 | `build.labels` | `Label=` [Build] | | 5.2.0 |
| 173 | `build.network` | `Network=` [Build] | | 5.2.0 |
| 174 | `build.no_cache` | `PodmanArgs=--no-cache` [Build] | | 5.2.0 |
| 175 | `build.secrets` | `Secret=` [Build] | | 5.2.0 |
| 176 | `build.tags` | `ImageTag=` [Build] | | 5.2.0 |

File: `mapper/network.go` — `.network` top-level unit

| # | Compose field | → Quadlet | Notes | Since |
|---|---|---|---|---|
| 177 | `networks.<name>.driver` | `Driver=` [Network] | | 4.4.0 |
| 178 | `networks.<name>.driver_opts` | `Options=` [Network] | map→key=value | 4.4.0 |
| 179 | `networks.<name>.ipam.driver` | `IPAMDriver=` [Network] | | 4.4.0 |
| 180 | `networks.<name>.ipam.config[].subnet` | `Subnet=` [Network] | multi-ipam | 4.4.0 |
| 181 | `networks.<name>.ipam.config[].gateway` | `Gateway=` [Network] | | 4.4.0 |
| 182 | `networks.<name>.ipam.config[].ip_range` | `IPRange=` [Network] | | 4.4.0 |
| 183 | `networks.<name>.internal` | `Internal=` [Network] | | 4.4.0 |
| 184 | `networks.<name>.enable_ipv6` | `IPv6=` [Network] | | 4.4.0 |
| 185 | `networks.<name>.external` | *(skip generation)* | no unit emitted | — |
| 186 | `networks.<name>.labels` | `Label=` [Network] | | 5.6.0 |
| 187 | `networks.<name>.dns` | `DNS=` [Network] | | 4.7.0 |
| 188 | `networks.<name>.interface_name` | `InterfaceName=` [Network] | | 5.6.0 |
| 189 | `networks.<name>.disable_dns` | `DisableDNS=` [Network] | | ? |
| 190 | `networks.<name>.delete_on_stop` | `NetworkDeleteOnStop=` [Network] | | 5.5.0 |

File: `mapper/volume.go` — `.volume` top-level unit

| # | Compose field | → Quadlet | Notes | Since |
|---|---|---|---|---|
| 191 | `volumes.<name>.driver` | `Driver=` [Volume] | | 4.7.0 |
| 192 | `volumes.<name>.driver_opts` | `Options=` [Volume] | | 6.0.0 |
| 193 | `volumes.<name>.external` | *(skip generation)* | no unit emitted | — |
| 194 | `volumes.<name>.labels` | `Label=` [Volume] | | ? |
| 195 | `volumes.<name>.name` | `VolumeName=` [Volume] | | 4.7.0 |
| 196 | `volumes.<name>.uid` | `UID=` [Volume] | | 6.0.0 |
| 197 | `volumes.<name>.gid` | `GID=` [Volume] | | 6.0.0 |
| 198 | `volumes.<name>.copy` | `Copy=` [Volume] | | ? |
| 199 | `volumes.<name>.device` | `Device=` [Volume] | | ? |
| 200 | `volumes.<name>.type` | `Type=` [Volume] | | ? |

---

## P4 — Unsupported (emit warnings only)

No mapping possible. Produce `WarningSkipped` for informational purposes.

| # | Compose field | Reason |
|---|---|---|
| 201 | `deploy.mode` | Swarm orchestration |
| 202 | `deploy.replicas` | Swarm orchestration |
| 203 | `deploy.placement.constraints` | Swarm orchestration |
| 204 | `deploy.placement.preferences` | Swarm orchestration |
| 205 | `deploy.endpoint_mode` | Swarm orchestration |
| 206 | `deploy.labels` | Swarm orchestration |
| 207 | `deploy.update_config` | Swarm orchestration |
| 208 | `deploy.rollback_config` | Swarm orchestration |
| 209 | `deploy.resources.reservations.devices` | Swarm orchestration |
| 210 | `extends` | handled by compose-go loader |
| 211 | `external_links` | legacy Docker |
| 212 | `links` | legacy Docker |
| 213 | `profiles` | runtime selection (comquad) |
| 214 | `scale` | replaces deploy.replicas |
| 215 | `domainname` | legacy Swarm |
| 216 | `credential_spec` | Windows-only |
| 217 | `isolation` | Windows/Swarm |
| 218 | `attach` | Docker CLI concept |
| 219 | `develop` | dev tooling |
| 220 | `volumes_from` | Docker-only |
| 221 | `volumes.subpath` | compose 2.27+ Docker engine only |
| 222 | `volumes.consistency` | inconsequential in podman |
| 223 | `cpu_count` | Windows/macOS |
| 224 | `cpu_percent` | Windows/macOS |
| 225 | `gpus` | no podman equivalent |
| 226 | `build.platforms` | multi-arch build (not applicable) |
| 227 | `build.extra_hosts` | no podman build equivalent |
| 228 | `networks.priority` | swarm/compose concept |
| 229 | `networks.driver_opts` | per-service network opts |
| 230 | `ipc: service:<name>` | no container IPC sharing |
| 231 | `pid: service:<name>` | no container PID namespace sharing |
| 232 | `networks.<name>.attachable` | swarm overlay only |
| 233 | `networks.<name>.ipam.config[].aux_addresses` | no quadlet directive |

---

## Summary

| Tier | Count | Done | Description |
|---|---|---|---|
| T0 | 51 | 35 | Trivial 1:1 passthrough (11 ❌ no compose-go field, 2 ⏳ moved up) |
| T1 | 38 | 30 | Easy formatting / type conversion (5 ❌ no compose-go field) |
| T2 | 39 | 0 | Medium complexity, version gating |
| T3 | 20 | 0 | Hard (P3 PodmanArgs, complex parsing) |
| T4 | 7 | 0 | Very hard (secrets/configs interceptor) |
| T5 | 42 | 0 | Structural (separate quadlet units) |
| P4 | 33 | 0 | Unsupported (warn only) |
| **Total** | **230** | **65** | |
