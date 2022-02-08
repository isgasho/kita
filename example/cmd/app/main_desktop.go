//go:build !ios && !android
// +build !ios,!android

package main

import (
	"fmt"
	"os"

	"github.com/uiez/uikit/app"
	"github.com/uiez/uikit/core/coregeom"
	"github.com/uiez/uikit/core/corelog"
	"github.com/uiez/uikit/example/pkg"
	"github.com/uiez/uikit/modules/envs"
	"github.com/uiez/uikit/modules/locales/i18n"
	"github.com/uiez/uikit/shell/platform"
)

func main() {
	profileStop := profileInit()
	a, err := app.New(
		"com.uiez.uikit.example",
	)
	if err != nil {
		corelog.Tag("App").Error("create app failed", err)
		return
	}
	if lang := os.Getenv("UI_LANG"); lang != "" {
		a.SetLanguage(lang)
	}

	var size coregeom.Point
	if s := envs.Get("WINDOW_SIZE"); s != "" {
		fmt.Sscanf(s, "%fx%f", &size.X, &size.Y)
	}
	titlebarStyle := envs.Get("WINDOW_TITLEBAR")
	if titlebarStyle == "" {
		titlebarStyle = platform.TitlebarStyleNormal
	}
	title := "UIKit Example"
	if t := os.Getenv("WINDOW_TITLE"); t != "" {
		title = t
	}
	resizable := envs.GetBool("WINDOW_RESIZABLE", true)
	borderless := envs.GetBool("WINDOW_BORDERLESS")

	_ = profileStop
	a.OnDestroy(func() {
		profileStop.SafeCancel()
	})
	a.Run(func() {
		_, err := a.CreateWindow(
			pkg.Content(),
			app.OptWindowTitle(i18n.Str(title)),
			app.OptWindowViewportSize(size),
			app.OptWindowTitleBarStyle(titlebarStyle),
			app.OptWindowResizable(resizable),
			app.OptWindowBorderless(borderless),
		)
		if err != nil {
			corelog.Tag("App").Error("create main window failed:", err)
			a.RequireClose()
			return
		}
	})
}
