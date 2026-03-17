package photon

import (
	"fmt"

	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/sqweek/dialog"
)

func Ask(app sandbox.Manifest, key string, permission sandbox.Permission) bool {
	var reason = permission.Reason

	if len(permission.Reason) == 0 {
		switch permission.Type {
		case sandbox.PermissionTypeRead:
			reason = "read the path '%s'"
		case sandbox.PermissionTypeReadWrite:
			reason = "access the path '%s'"
		}

		reason = fmt.Sprintf(reason, key)
	}

	return dialog.Message("%s wants to %s. Do you want to give the app permission?", app.Application.Name, reason).Title("Are you sure?").YesNo()
}
