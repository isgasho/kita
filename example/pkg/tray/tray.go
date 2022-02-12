package tray

import (
	"fmt"
	"image"
	"math/rand"
	"strings"
	"time"

	"github.com/uiez/uikit/app"
	"github.com/uiez/uikit/components/filechoose"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/coreevent"
	"github.com/uiez/uikit/core/coregeom"
	"github.com/uiez/uikit/core/coregfx"
	"github.com/uiez/uikit/core/corert"
	"github.com/uiez/uikit/core/coreui"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
	"github.com/uiez/uikit/modules/gfx/colors"
	"github.com/uiez/uikit/modules/gfx/images/imgcodec"
	"github.com/uiez/uikit/modules/locales/i18n"
	"github.com/uiez/uikit/modules/services"
	"github.com/uiez/uikit/shell/platform"
)

type appState struct {
	corewidget.StateBase

	app         *app.App
	window      *app.Window
	trayService corert.Optional[services.Tray]

	action   string
	errorMsg string

	trayColor    string
	trayIconPath string

	trayItem   services.TrayItem
	trayWindow *app.Window

	trayIconChoose *filechoose.State
}

func newAppState(node corewidget.ComponentNode) *appState {
	s := &appState{}

	coreui.GetModuleInto(node.View(), &s.app)
	coreui.GetModuleInto(node.View(), &s.window)
	coreui.GetOptionalModuleInto(node.View(), &s.trayService)

	rand.Seed(time.Now().UnixNano())
	s.trayIconChoose = filechoose.NewState(node.View())
	return s
}

func (s *appState) Destroy() {
	s.StateBase.Destroy()
	s.trayIconChoose.Destroy()
	s.trayDestroy()
}

func (s *appState) onAction(action string) {
	s.errorMsg = ""
	s.action = action

	s.Notify()
}

func (s *appState) onError(errMsg string) {
	s.action = ""
	s.errorMsg = errMsg
	s.Notify()
}

func (s *appState) newColorTrayIcon(color coregfx.Color) *image.RGBA {
	rgba := image.NewRGBA(image.Rectangle{
		Max: image.Point{X: 32, Y: 32},
	})
	size := rgba.Bounds().Size()
	for i := 0; i < size.X; i++ {
		for j := 0; j < size.Y; j++ {
			rgba.Set(i, j, color)
		}
	}
	return rgba
}

func (s *appState) randTrayColor() {
	colornames := []string{"red", "orange", "white", "yellow", "green", "fuchsia", "purple", "black", "blue", "cyan"}
	name := colornames[rand.Intn(len(colornames))]

	s.trayIconPath = ""
	s.trayColor = name
	if s.trayItem != nil {
		s.syncTrayIcon()
	}
	s.Notify()
}

func (s *appState) chooseTrayIcon() {
	exts := []services.FileFilter{
		{Description: "Images", MIMEs: []string{"image/png", "image/jpeg"}},
	}
	s.trayIconChoose.Open(i18n.Str("choose tray icon"), exts, services.FileDialogAllowFile, func(selected bool, paths []string) {
		if selected && len(paths) > 0 {
			s.trayIconPath = paths[0]
			s.trayColor = ""
			if s.trayItem != nil {
				s.syncTrayIcon()
			}
			s.Notify()
		}
	})
}

func (s *appState) syncTrayIcon() {
	if s.trayIconPath != "" {
		icon, err := imgcodec.DecodeFileRGBA(s.trayIconPath)
		if err != nil {
			s.onError(fmt.Sprintf("decode image failed: %s", err.Error()))
		} else {
			s.trayItem.SetIcon(icon)
		}
	} else if s.trayColor != "" {
		s.trayItem.SetIcon(s.newColorTrayIcon(coregfx.NewColor(colors.ColorName(s.trayColor))))
	} else {
		s.trayItem.SetIcon(s.newColorTrayIcon(colors.Orange))
	}
}

func (s *appState) initTray() {
	s.syncTrayIcon()
	s.trayItem.SetTooltip("tray icon")
	s.trayItem.SetAction("click action", func() {
		if s.trayWindow == nil {
			w, err := s.app.CreateWindow(
				_TrayWindow(),
				app.OptWindowViewportSize(coregeom.Point{X: 300, Y: 400}),
				app.OptWindowResizable(false),
				app.OptWindowIsPopup(true),
				app.OptWindowTitleBarStyle(platform.TitlebarStyleHide),
				app.OptWindowPosition(app.WinposTrayPopup(s.trayItem.Bounds())),
			)
			if err != nil {
				s.onError(fmt.Sprintf("create tray window failed: %s", err.Error()))
			} else {
				s.trayWindow = w
				s.trayWindow.OnDestroy(func() {
					s.trayWindow = nil
				})
			}
		} else {
			s.trayWindow.SetVisible(!s.trayWindow.State().Visible)
		}
	}, services.MenuItems{
		{Label: "close window", Action: "window.close"},
		{Label: "test action", Action: "test"},
		{Label: "test action2", Action: "test2"},
		{Label: "remove tray", Action: "tray.remove"},
		{Label: "app quit", Action: "app.quit"},
	}, func(action string) {
		s.onAction(action)

		switch action {
		case "window.close":
			s.window.RequestClose()
		case "app.quit":
			s.app.RequireClose()
		case "tray.remove":
			s.toggleTray()
		}
	})
}

func (s *appState) toggleTray() {
	if !s.trayService.OK() {
		s.onError("tray service not available")
		return
	}

	if s.trayItem == nil {
		var err error
		s.trayItem, err = s.trayService.Value().CreateItem()
		if err != nil {
			s.onError("add tray item failed:" + err.Error())
			return
		}
		s.initTray()
	} else {
		s.trayDestroy()
	}
	s.Notify()
}

func (s *appState) trayDestroy() {
	if s.trayWindow != nil {
		s.trayWindow.Destroy()
		s.trayWindow = nil
	}
	if s.trayItem != nil {
		s.trayItem.Destroy()
		s.trayItem = nil
	}
}

var _TrayWindow = corewidget.Component("TrayWindow", func(node corewidget.ComponentNode) corewidget.Widget {
	focusChange, currState := corewidget.UseSettingsChange(node, func(old, curr coreui.WindowState) bool {
		return old.Focused != curr.Focused
	})
	if focusChange && !currState.Focused {
		w, _ := coreui.GetModule[*app.Window](node.View())
		w.SetVisible(false)
	}

	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("window.scss"))
	return dom.Div(
		dom.Stylesheets(stylesheets),
		dom.OnWindowShortcut(func() []coreevent.Shortcut {
			return []coreevent.Shortcut{{Key: coreevent.KeyEscape}}
		}, func(ctx coreevent.PropagateContext, shortcut coreevent.Shortcut) {
			w, _ := coreui.GetModule[*app.Window](node.View())
			w.SetVisible(false)
		}),
	).Children(dom.Span(dom.SpanText("popup window here")))
})

var App = corewidget.Component("Tray", func(node corewidget.ComponentNode) corewidget.Widget {
	as := corewidget.UseBuilder(node, newAppState)
	var text []string
	if as.errorMsg != "" {
		text = append(text, "error: "+as.errorMsg)
	} else if as.action == "" {
		text = append(text, "action: no action")
	} else {
		text = append(text, "action:"+as.action)
	}
	if as.trayColor == "" {
		text = append(text, "icon color: not randed")
	} else {
		text = append(text, "color: "+as.trayColor)
	}
	if as.trayIconPath == "" {
		text = append(text, "icon file: not selected")
	} else {
		text = append(text, "icon file: "+as.trayIconPath)
	}

	enableOrDisable := func(enabled bool) string {
		if enabled {
			return "disable"
		}
		return "enable"
	}
	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("tray.scss"))
	return dom.Div(dom.Stylesheets(stylesheets)).Children(
		dom.Span(dom.SpanText(strings.Join(text, "\n"))),

		dom.Button(
			dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
				as.randTrayColor()
			}),
		).Child(dom.Label(dom.LabelText("rand tray color"))),
		dom.Button(
			dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
				as.chooseTrayIcon()
			}),
		).Child(dom.Label(dom.LabelText("choose tray icon"))),

		dom.Button(
			dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
				as.toggleTray()
			}),
		).Child(dom.Label(dom.LabelText(enableOrDisable(as.trayItem != nil)+" tray"))),
	)
})
