//go:build !debug
// +build !debug

package main

import "github.com/uiez/uikit/core/corebase"

func profileInit() corebase.CancelFunc {
	return func() {}
}
