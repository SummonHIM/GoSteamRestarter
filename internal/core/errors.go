package core

import "errors"

var (
	ErrSteamNotFound    = errors.New("steam not found")
	ErrSteamNotRunning  = errors.New("steam not running")
	ErrPermissionDenied = errors.New("permission denied")
)
