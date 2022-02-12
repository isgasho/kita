package window

import (
	"fmt"
	"strconv"

	"github.com/uiez/uikit/app"
	"github.com/uiez/uikit/components/menu"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/coreevent"
	"github.com/uiez/uikit/core/coregeom"
	"github.com/uiez/uikit/core/corelog"
	"github.com/uiez/uikit/core/coreui"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/dom/styles"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
	"github.com/uiez/uikit/modules/gfx/colors"
	"github.com/uiez/uikit/modules/locales/i18n"
	"github.com/uiez/uikit/modules/services"
	"github.com/uiez/uikit/modules/texts"
	"github.com/uiez/uikit/shell/platform"
)

var transparentPanel = corewidget.PComponent("TransparentPanel", func(node corewidget.ComponentNode, initialTransparency float64) corewidget.Widget {
	transparency := corewidget.UseValue(node, initialTransparency)
	transparencyVal := corebase.Clamp(transparency.Get(), 0, 1)
	return dom.Div(
		dom.Class("panel"),
		dom.InlineStyle(
			styles.PropBackgroundColor.NewRuleProp(colors.Opacity(colors.Gray, transparencyVal)),
		),
	).Children(
		dom.Div(
			dom.Class("form"),
		).Children(
			dom.Label(dom.LabelText("Opacity:")),

			dom.Input(
				dom.InputPlaceholder("input opacity value"),
				dom.InputOnChange(func(text texts.Text) {
					v, _ := strconv.ParseFloat(texts.CopyTextDataAsString(text), 32)
					transparency.Set(v)
				}),
				dom.InputValue(texts.Runes(fmt.Sprintf("%v", transparency.Get()))),
			),
		),
	)
})

var transparentWindow = corewidget.Component("TransparentWindow", func(node corewidget.ComponentNode) corewidget.Widget {
	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("window.scss"))
	return dom.Div(
		dom.Stylesheets(stylesheets),
		dom.Id("transparent"),
	).Children(
		transparentPanel(0.3),
		transparentPanel(0.7),
	)
})

func createTransparentWindow(application *app.App) {
	w, err := application.CreateWindow(
		transparentWindow(),
		app.OptWindowTitle(i18n.Str("Transparent window")),
		app.OptWindowViewportSize(coregeom.Point{X: 600, Y: 400}),
		app.OptWindowTransparent(true),
	)
	if err != nil {
		corelog.Tag("Window").Error("window create failed:", err)
		return
	}
	w.OnFocusChange(func(focused bool) {
		if !focused {
			w.RequestClose()
		}
	})
}

var dynamicSizeWindow = corewidget.Component("DynamicSizeWindow", func(node corewidget.ComponentNode) corewidget.Widget {
	number := corewidget.UseValue[texts.Text](node, texts.Runes{})

	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("window.scss"))
	return dom.Div(dom.Id("dynamic-size"), dom.Stylesheets(stylesheets)).Children(
		dom.Input(
			dom.Autofocus(true),
			dom.InputType(services.TextInputTypeNumber),
			dom.InputValue(number.Get()),
			dom.InputOnChange(number.Set),
			dom.InputPlaceholder("type a number, show items with count N%30"),
		),
		dom.Div(dom.Id("dropdown")).Compose(func(emit corewidget.ComposeEmitFunc) {
			n, _ := strconv.Atoi(string(texts.CopyTextDataAsString(number.Get())))
			for i := 0; i < n%30; i++ {
				emit(dom.Label(dom.Class("dropdown-item"), dom.LabelText(strconv.Itoa(i+n))))
			}
		}),
	)
})

func createDynamicSizeWindow(application *app.App) {
	w, err := application.CreateWindow(
		dynamicSizeWindow(),
		app.OptWindowViewportMatchContentSize(true),
		app.OptWindowTitleBarStyle(platform.TitlebarStyleHide),
		app.OptWindowIsPopup(true),
		app.OptWindowPosition(app.WinposAlignToScreen(coregeom.Alignment{X: 0.5, Y: 0.2})),
		app.OptWindowViewportSize(coregeom.Point{X: 800, Y: 60}),
	)
	if err != nil {
		corelog.Tag("Window").Error("window create failed:", err)
		return
	}
	w.OnFocusChange(func(focused bool) {
		if !focused {
			w.RequestClose()
		}
	})
}

var sizeLimitWindow = corewidget.Component("SizeLimitWindow", func(node corewidget.ComponentNode) corewidget.Widget {
	enable := corewidget.UseValue(node, false)

	var menubar services.MenuItems
	if enable.Get() {
		menubar = services.MenuItems{
			{
				Label: "test",
			},
		}
	}
	menu.UseWindowMenubar(node, menubar, func(action string) {})
	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("window.scss"))
	return dom.Div(dom.Id("size-limit"), dom.Stylesheets(stylesheets)).Children(
		dom.Button(dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
			enable.Set(!enable.Get())
		})).Child(dom.Label(dom.LabelText("toggle menubar"))),
		dom.Label(
			dom.LabelText("size: 200 ~ 400, default 300"),
		),
	)
})

func createSizeLimitedWindow(application *app.App) {
	w, err := application.CreateWindow(
		sizeLimitWindow(),
		app.OptWindowViewportSize(coregeom.Point{X: 300, Y: 300}),
		app.OptWindowViewportMinSize(coregeom.Point{X: 200, Y: 200}),
		app.OptWindowViewportMaxSize(coregeom.Point{X: 400, Y: 400}),
	)
	if err != nil {
		corelog.Tag("Window").Error("window create failed:", err)
		return
	}
	w.OnFocusChange(func(focused bool) {
		if !focused {
			w.RequestClose()
		}
	})
}

var modalWindowContent = corewidget.Component("ModalWindow", func(node corewidget.ComponentNode) corewidget.Widget {
	onWindowKey := func(ctx coreevent.PropagateContext, ele coredom.Element, key coreevent.KeyEvent) {
		var window *app.Window
		coreui.GetModuleInto(ele.View(), &window)
		if key.Action == coreevent.KeyActionDown && key.Key == coreevent.KeyEscape {
			window.SetVisible(!window.State().Visible)
			window.Destroy()

			ctx.StopPropagation()
			ctx.Finish()
		}
	}

	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("window.scss"))
	return dom.Div(
		dom.Stylesheets(stylesheets),
		dom.Id("modal"),
		dom.OnWindowKey(onWindowKey),
	).Children(
		dom.Label(
			dom.LabelText("modal window"),
		),
	)
})

func createModalWindow(application *app.App, window *app.Window) {
	w, err := application.CreateWindow(
		modalWindowContent(),
		app.OptWindowParent(window),
		app.OptWindowIsModal(true),
		app.OptWindowViewportSize(coregeom.Point{X: 300, Y: 300}),
	)
	if err != nil {
		corelog.Tag("Window").Error("window create failed:", err)
		return
	}
	_ = w
}

var App = corewidget.Component("Window", func(node corewidget.ComponentNode) corewidget.Widget {
	var (
		application *app.App
		window      *app.Window
	)
	coreui.GetModuleInto(node.View(), &application)
	coreui.GetModuleInto(node.View(), &window)

	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("window.scss"))
	return dom.Div(dom.Id("pkg"), dom.Stylesheets(stylesheets)).Children(
		dom.Button(
			dom.OnTap(func(_ coredom.Element, _ recognizers.TapDetails) {
				createTransparentWindow(application)
			}),
		).Child(dom.Label(dom.LabelText("create transparent window"))),
		dom.Button(
			dom.OnTap(func(_ coredom.Element, _ recognizers.TapDetails) {
				createDynamicSizeWindow(application)
			}),
		).Child(dom.Label(dom.LabelText("create dynamic size window"))),
		dom.Button(
			dom.OnTap(func(_ coredom.Element, _ recognizers.TapDetails) {
				createSizeLimitedWindow(application)
			}),
		).Child(dom.Label(dom.LabelText("create size limited window"))),
		dom.Button(
			dom.OnTap(func(_ coredom.Element, _ recognizers.TapDetails) {
				createModalWindow(application, window)
			}),
		).Child(dom.Label(dom.LabelText("create modal window"))),
	)
})
