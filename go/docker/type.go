package docker

import (
	"time"
)

// ContainerConfig reference for config.json in docker runtime
type ContainerConfig struct {
	State struct {
		Running           bool        `json:"Running"`
		Paused            bool        `json:"Paused"`
		Restarting        bool        `json:"Restarting"`
		OOMKilled         bool        `json:"OOMKilled"`
		RemovalInProgress bool        `json:"RemovalInProgress"`
		Dead              bool        `json:"Dead"`
		Pid               int         `json:"Pid"`
		ExitCode          int         `json:"ExitCode"`
		Error             string      `json:"Error"`
		StartedAt         time.Time   `json:"StartedAt"`
		FinishedAt        time.Time   `json:"FinishedAt"`
		Health            interface{} `json:"Health"`
	} `json:"State"`
	ID      string    `json:"ID"`
	Created time.Time `json:"Created"`
	Managed bool      `json:"Managed"`
	Path    string    `json:"Path"`
	Args    []string  `json:"Args"`
	Config  struct {
		Hostname     string              `json:"Hostname"`
		Domainname   string              `json:"Domainname"`
		User         string              `json:"User"`
		AttachStdin  bool                `json:"AttachStdin"`
		AttachStdout bool                `json:"AttachStdout"`
		AttachStderr bool                `json:"AttachStderr"`
		Tty          bool                `json:"Tty"`
		OpenStdin    bool                `json:"OpenStdin"`
		StdinOnce    bool                `json:"StdinOnce"`
		Env          []string            `json:"Env"`
		Cmd          []string            `json:"Cmd"`
		Image        string              `json:"Image"`
		Volumes      map[string]struct{} `json:"Volumes"`
		WorkingDir   string              `json:"WorkingDir"`
		Entrypoint   []string            `json:"Entrypoint"`
		OnBuild      []string            `json:"OnBuild"`
		Labels       map[string]string   `json:"Labels"`
	} `json:"Config"`
	Image           string `json:"Image"`
	NetworkSettings struct {
		Bridge                 string `json:"Bridge"`
		SandboxID              string `json:"SandboxID"`
		HairpinMode            bool   `json:"HairpinMode"`
		LinkLocalIPv6Address   string `json:"LinkLocalIPv6Address"`
		LinkLocalIPv6PrefixLen int    `json:"LinkLocalIPv6PrefixLen"`
		Networks               struct {
			Bridge struct {
				IPAMConfig          interface{} `json:"IPAMConfig"`
				Links               interface{} `json:"Links"`
				Aliases             interface{} `json:"Aliases"`
				NetworkID           string      `json:"NetworkID"`
				EndpointID          string      `json:"EndpointID"`
				Gateway             string      `json:"Gateway"`
				IPAddress           string      `json:"IPAddress"`
				IPPrefixLen         int         `json:"IPPrefixLen"`
				IPv6Gateway         string      `json:"IPv6Gateway"`
				GlobalIPv6Address   string      `json:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int         `json:"GlobalIPv6PrefixLen"`
				MacAddress          string      `json:"MacAddress"`
				DriverOpts          interface{} `json:"DriverOpts"`
				IPAMOperational     bool        `json:"IPAMOperational"`
			} `json:"bridge"`
		} `json:"Networks"`
		Service interface{} `json:"Service"`
		Ports   struct {
		} `json:"Ports"`
		SandboxKey             string      `json:"SandboxKey"`
		SecondaryIPAddresses   interface{} `json:"SecondaryIPAddresses"`
		SecondaryIPv6Addresses interface{} `json:"SecondaryIPv6Addresses"`
		IsAnonymousEndpoint    bool        `json:"IsAnonymousEndpoint"`
		HasSwarmEndpoint       bool        `json:"HasSwarmEndpoint"`
	} `json:"NetworkSettings"`
	LogPath                string `json:"LogPath"`
	Name                   string `json:"Name"`
	Driver                 string `json:"Driver"`
	Os                     string `json:"OS"`
	MountLabel             string `json:"MountLabel"`
	ProcessLabel           string `json:"ProcessLabel"`
	RestartCount           int    `json:"RestartCount"`
	HasBeenStartedBefore   bool   `json:"HasBeenStartedBefore"`
	HasBeenManuallyStopped bool   `json:"HasBeenManuallyStopped"`
	MountPoints            map[string]struct {
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		Rw          bool   `json:"RW"`
		Name        string `json:"Name"`
		Driver      string `json:"Driver"`
		Type        string `json:"Type"`
		Propagation string `json:"Propagation"`
		Spec        struct {
			Type   string `json:"Type"`
			Source string `json:"Source"`
			Target string `json:"Target"`
		} `json:"Spec"`
		SkipMountpointCreation bool `json:"SkipMountpointCreation"`
	} `json:"MountPoints"`
	SecretReferences  interface{} `json:"SecretReferences"`
	ConfigReferences  interface{} `json:"ConfigReferences"`
	AppArmorProfile   string      `json:"AppArmorProfile"`
	HostnamePath      string      `json:"HostnamePath"`
	HostsPath         string      `json:"HostsPath"`
	ShmPath           string      `json:"ShmPath"`
	ResolvConfPath    string      `json:"ResolvConfPath"`
	SeccompProfile    string      `json:"SeccompProfile"`
	NoNewPrivileges   bool        `json:"NoNewPrivileges"`
	LocalLogCacheMeta struct {
		HaveNotifyEnabled bool `json:"HaveNotifyEnabled"`
	} `json:"LocalLogCacheMeta"`
}
