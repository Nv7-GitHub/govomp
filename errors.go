package govomp

import "errors"

var (
	ErrNoSuitableMemoryTypeFound = errors.New("govomp: no suitable memory type found")
	ErrFailedToCopyMemory        = errors.New("govomp: failed to copy memory")
)
