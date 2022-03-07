package utils

import "github.com/buliqioqiolibusdo/demp-core/interfaces"

func GetUserFromArgs(args ...interface{}) (u interfaces.User) {
	for _, arg := range args {
		switch arg.(type) {
		case interfaces.User:
			var ok bool
			u, ok = arg.(interfaces.User)
			if ok {
				return u
			}
		}
	}
	return nil
}
