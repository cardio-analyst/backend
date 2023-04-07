package model

import "errors"

var ErrIPIsNotInWhitelist = errors.New("user ip is not in the whitelist")
