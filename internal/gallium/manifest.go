package gallium

import (
	"regexp"

	"github.com/arvasyn/warden/internal/pkg/apperr"
)

type Manifest struct {
	Version int      `yaml:"version"`
	Type    UnitType `yaml:"type"`

	// Application
	Application Application `yaml:"application"`

	// Unit is only used if type is set to service
	Unit *Unit `yaml:"unit"`

	// This option is *only* available to services.
	// Applications *must* be sandboxed. Services must be code signed in order to use this option.
	UseSandbox bool    `yaml:"use_sandbox,omitempty"`
	Sandbox    Sandbox `yaml:"sandbox"`
}

func (m *Manifest) Validate() error {
	bundle, err := regexp.Compile("^[A-Za-z0-9_]+(\\.[A-Za-z0-9_]+)+$")
	if err != nil {
		return err
	}

	if !bundle.MatchString(m.Application.Bundle) {
		return apperr.ErrInvalidApplicationBundle
	}

	if m.Type == UnitTypeApplication && !m.UseSandbox {
		return apperr.ErrApplicationNoDisableSandbox
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
	Namespaces   Namespaces            `yaml:"namespaces,omitempty"`
	Filesystem   Filesystem            `yaml:"filesystem,omitempty"`
	Env          EnvPolicy             `yaml:"env,omitempty"`
	Capabilities Capabilities          `yaml:"capabilities,omitempty"`
	Resources    Resources             `yaml:"resources,omitempty"`
	Portals      []Portals             `yaml:"portals,omitempty"`
	Permissions  map[string]Permission `yaml:"permissions,omitempty"`
}

type Portals string

const (
	PortalAccount         Portals = "org.freedesktop.portal.Account"
	PortalAppChooser      Portals = "org.freedesktop.portal.AppChooser"
	PortalBackground      Portals = "org.freedesktop.portal.Background"
	PortalCamera          Portals = "org.freedesktop.portal.Camera"
	PortalClipboard       Portals = "org.freedesktop.portal.Clipboard"
	PortalDevice          Portals = "org.freedesktop.portal.Device"
	PortalDynamicLauncher Portals = "org.freedesktop.portal.DynamicLauncher"
	PortalEmail           Portals = "org.freedesktop.portal.Email"
	PortalFileChooser     Portals = "org.freedesktop.portal.FileChooser"
	PortalGlobalShortcuts Portals = "org.freedesktop.portal.GlobalShortcuts"
	PortalInhibit         Portals = "org.freedesktop.portal.Inhibit"
	PortalInputCapture    Portals = "org.freedesktop.portal.InputCapture"
	PortalLocation        Portals = "org.freedesktop.portal.Location"
	PortalLockdown        Portals = "org.freedesktop.portal.Lockdown"
	PortalMemoryMonitor   Portals = "org.freedesktop.portal.MemoryMonitor"
	PortalNetworkMonitor  Portals = "org.freedesktop.portal.NetworkMonitor"
	PortalNotification    Portals = "org.freedesktop.portal.Notification"
	PortalOpenURI         Portals = "org.freedesktop.portal.OpenURI"
	PortalPrint           Portals = "org.freedesktop.portal.Print"
	PortalProxyResolver   Portals = "org.freedesktop.portal.ProxyResolver"
	PortalRemoteDesktop   Portals = "org.freedesktop.portal.RemoteDesktop"
	PortalRequest         Portals = "org.freedesktop.portal.Request"
	PortalScreenCast      Portals = "org.freedesktop.portal.ScreenCast"
	PortalSecret          Portals = "org.freedesktop.portal.Secret"
	PortalSettings        Portals = "org.freedesktop.portal.Settings"
	PortalScreenshot      Portals = "org.freedesktop.portal.Screenshot"
	PortalSession         Portals = "org.freedesktop.portal.Session"
	PortalTrash           Portals = "org.freedesktop.portal.Trash"
	PortalWallpaper       Portals = "org.freedesktop.portal.Wallpaper"
)

type Namespaces map[string]bool

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

type Resources struct {
	MemoryMB       int `yaml:"memory_mb,omitempty"`
	CPUPercent     int `yaml:"cpu_percent,omitempty"`
	TimeoutSeconds int `yaml:"timeout_seconds,omitempty"`
	MaxProcesses   int `yaml:"max_processes,omitempty"`
}

type Permission struct {
	Type   PermissionType `yaml:"type"`
	Reason string         `yaml:"reason,omitempty"`
}

type PermissionType string

const (
	PermissionTypeRead      PermissionType = "READ"
	PermissionTypeReadWrite PermissionType = "READ+WRITE"
)
