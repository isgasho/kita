package main

import (
	"github.com/uiez/uikit/app"
	"github.com/uiez/uikit/core/corelog"
	"github.com/uiez/uikit/example/pkg"
)

func main() {
	a, err := app.New("com.uiez.uikit.example")
	if err != nil {
		corelog.Tag("App").Error("create app failed", err)
		return
	}

	content := pkg.Content()
	a.Run(func() {
		_, err := a.CreateWindow(
			content,
		)
		if err != nil {
			corelog.Tag("App").Error("create main window failed:", err)
			a.RequireClose()
			return
		}
	})
}
