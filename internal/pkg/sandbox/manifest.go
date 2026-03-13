package sandbox

import "errors"

type Manifest struct {
	Version int      `yaml:"version"`
	Type    UnitType `yaml:"type"`

	// Application
	Application Application `yaml:"application"`

	// Unit is only used if type is set to service
	Unit *Unit `yaml:"unit"`

	// This option is *only* available to services.
	// Applications *must* be sandboxed. Services must be code signed in order to use this option.
	UseSandbox bool     `yaml:"use_sandbox"`
	Sandbox    *Sandbox `yaml:"sandbox"`
}

func (m *Manifest) Validate() error {
	if m.Type == UnitTypeApplication && !m.UseSandbox {
		return errors.New("applications must have use_sandbox set to true")
	}
	return nil
}

type Application struct {
	Name        string   `yaml:"name"`
	Bundle      string   `yaml:"bundle"`
	Description string   `yaml:"description,omitempty"`
	Exec        string   `yaml:"exec"`
	Args        []string `yaml:"args,omitempty"`
	// If set to true, user specified arguments will be added after the already specified ones
	AppendUserArgs bool `yaml:"append_user_args,omitempty"`
}

type Unit struct {
	User        string      `yaml:"user,omitempty"`
	Group       string      `yaml:"group,omitempty"`
	Supervision Supervision `yaml:"supervision,omitempty"`
}

type UnitType string

const (
	UnitTypeService     UnitType = "service"
	UnitTypeApplication UnitType = "application"
)

type Supervision struct {
	Restart      string `yaml:"restart,omitempty"`
	RestartDelay int    `yaml:"restart_delay,omitempty"`
	StopSignal   string `yaml:"stop_signal,omitempty"`
}

type Sandbox struct {
	// All namespaces will default to true if not explicitly declared
	Namespaces   Namespaces    `yaml:"namespaces,omitempty"`
	Filesystem   Filesystem    `yaml:"filesystem,omitempty"`
	Env          EnvPolicy     `yaml:"env,omitempty"`
	Capabilities Capabilities  `yaml:"capabilities,omitempty"`
	Process      ProcessPolicy `yaml:"process,omitempty"`
	Resources    Resources     `yaml:"resources,omitempty"`
	Portals      []PortalName  `yaml:"portals,omitempty"`
}

type Namespaces struct {
	Net  *bool `yaml:"net,omitempty"`
	Pid  *bool `yaml:"pid,omitempty"`
	IPC  *bool `yaml:"ipc,omitempty"`
	UTS  *bool `yaml:"uts,omitempty"`
	User *bool `yaml:"user,omitempty"`
}

type Filesystem struct {
	Mounts []Mount `yaml:"mounts,omitempty"`
}

type Mount struct {
	Type   MountType `yaml:"type"`
	Source string    `yaml:"source,omitempty"`
	Target string    `yaml:"target"`
}

type MountType string

const (
	// MountTypeBind --bind (read-write)
	MountTypeBind MountType = "bind"

	// MountTypeROBind --ro-bind (read only)
	MountTypeROBind MountType = "ro_bind"

	// MountTypeTmpfs --tmpfs
	MountTypeTmpfs MountType = "tmpfs"

	// MountTypeProc --proc
	MountTypeProc MountType = "proc"

	// MountTypeDev --dev
	MountTypeDev MountType = "dev"
)

type EnvPolicy struct {
	// The environment variables to passthrough to the application (e.g. $HOME)
	Passthrough []string          `yaml:"passthrough,omitempty"`
	Set         map[string]string `yaml:"set,omitempty"`
	Unset       []string          `yaml:"unset,omitempty"`
}

type Capabilities struct {
	DropAll bool     `yaml:"drop_all,omitempty"`
	Add     []string `yaml:"add,omitempty"`
}

type ProcessPolicy struct {
	DieWithParent bool `yaml:"die_with_parent,omitempty"`
	NewSession    bool `yaml:"new_session,omitempty"`
}

type Resources struct {
	MemoryMB       int `yaml:"memory_mb,omitempty"`
	CPUPercent     int `yaml:"cpu_percent,omitempty"`
	TimeoutSeconds int `yaml:"timeout_seconds,omitempty"`
	MaxProcesses   int `yaml:"max_processes,omitempty"`
}

type PortalName string

const (
	PortalFileChooser   PortalName = "file-chooser"
	PortalNetwork       PortalName = "network"
	PortalCamera        PortalName = "camera"
	PortalNotifications PortalName = "notifications"
	PortalSecret        PortalName = "secret"
	PortalOpenURI       PortalName = "open-uri"
)
