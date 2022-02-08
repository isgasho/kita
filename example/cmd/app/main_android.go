package main

import (
	"github.com/uiez/uikit/app"
	"github.com/uiez/uikit/core/corelog"
	"github.com/uiez/uikit/example/pkg"
	"github.com/uiez/uikit/shell/platform/android"
)

func main() {
	android.SetWindowContentBuilder(func(args string) (any, error) {
		return pkg.Content(), nil
	})

	a, err := app.New("com.uiez.uikit.example")
	if err != nil {
		corelog.Tag("App").Error("create app failed", err)
		return
	}
	// android app will automatically create first activity window
	a.Run(nil)
}
