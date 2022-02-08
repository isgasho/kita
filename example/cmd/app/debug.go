//go:build debug
// +build debug

package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"

	"github.com/pkg/profile"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/corelog"
)

func profileInit() corebase.CancelFunc {
	var profileStop corebase.CancelFunc
	switch os.Getenv("PROFILE") {
	case "cpu":
		profileStop = profile.Start(profile.CPUProfile).Stop
	case "mem":
		profileStop = profile.Start(profile.MemProfile).Stop
	case "http", "":
		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			corelog.Tag("Profile").Error("listen for pprof failed:", err)
		} else {
			var wg sync.WaitGroup
			profileStop = func() {
				lis.Close()
				wg.Wait()
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				http.Serve(lis, http.DefaultServeMux)
			}()
		}
	}
	return profileStop
}
