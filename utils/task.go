package utils

import "github.com/buliqioqiolibusdo/demp-core/constants"

func IsCancellable(status string) bool {
	switch status {
	case constants.TaskStatusPending,
		constants.TaskStatusRunning:
		return true
	default:
		return false
	}
}
