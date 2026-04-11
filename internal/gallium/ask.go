package gallium

import (
	"fmt"

	"github.com/sqweek/dialog"
)

func Ask(app Manifest, key string, permission Permission) bool {
	var reason = permission.Reason

	if len(permission.Reason) == 0 {
		switch permission.Type {
		case PermissionTypeRead:
			reason = "read the path '%v'"
		case PermissionTypeReadWrite:
			reason = "access the path '%v'"
		}

		reason = fmt.Sprintf(reason, key)
	}

	return dialog.Message(
		"%v wants to %v. Do you want to give the app permission?",
		app.Application.Name, reason,
	).
		Title("Are you sure?").
		YesNo()
}
