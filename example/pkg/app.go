package pkg

import (
	"strconv"

	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	_ "github.com/uiez/uikit/example/assets"
	"github.com/uiez/uikit/example/pkg/basic"
	"github.com/uiez/uikit/example/pkg/edit"
	"github.com/uiez/uikit/example/pkg/filebrowser"
	"github.com/uiez/uikit/example/pkg/hotkey"
	"github.com/uiez/uikit/example/pkg/i18n"
	"github.com/uiez/uikit/example/pkg/menu"
	"github.com/uiez/uikit/example/pkg/popup"
	"github.com/uiez/uikit/example/pkg/router"
	"github.com/uiez/uikit/example/pkg/tray"
	"github.com/uiez/uikit/example/pkg/window"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
)

//#uikit asset: *.scss */*.scss

type appModule struct {
	Title   string
	Content corewidget.Widget
}

func (m *appModule) key(idx int) string {
	return m.Title + "-" + strconv.Itoa(idx)
}

type appState struct {
	corewidget.StateBase

	modules    []appModule
	currModule string
}

func newAppState(node corewidget.ComponentNode, modules []appModule) *appState {
	var initialModule string
	if len(modules) > 0 {
		initialModule = modules[0].key(0)
	}
	return &appState{
		modules:    modules,
		currModule: initialModule,
	}
}
func (s *appState) Update(oldProps, props []appModule) {}
func (s *appState) switchModule(mod string) {
	if s.currModule != mod {
		s.currModule = mod
		s.Notify()
	}
}

var sidebar = corewidget.Component("Sidebar", func(node corewidget.ComponentNode) corewidget.Widget {
	state := corewidget.UsePBuilderConsumer(node, newAppState).Value()

	return dom.Div(
		dom.Id("sidebar"),
	).Compose(func(emit corewidget.ComposeEmitFunc) {
		for i, mod := range state.modules {
			key := mod.key(i)
			if i > 0 {
				dom.Spacer().Emit(emit)
			}
			dom.Label.Key(key)(
				dom.Class("sidebar-item", dom.OptionClass(key == state.currModule, "active")),
				dom.LabelText(mod.Title),
				dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
					state.switchModule(key)
				}),
			).Emit(emit)
		}
	})
})

var content = corewidget.Component("Content", func(node corewidget.ComponentNode) corewidget.Widget {
	state := corewidget.UsePBuilderConsumer(node, newAppState).Value()
	currMod := state.currModule

	var active appModule
	var key string
	for i := range state.modules {
		mod := state.modules[i]
		k := mod.key(i)
		if k == currMod {
			active = mod
			key = k
		}
	}

	return dom.Div(
		dom.Id("content"),
	).Compose(func(emit corewidget.ComposeEmitFunc) {
		if active.Content.Class != nil {
			dom.Div.Key(key)().
				Children(active.Content).Emit(emit)
		}
	})
})

var _App = corewidget.PComponent("App", func(node corewidget.ComponentNode, modules []appModule) corewidget.Widget {
	var children []corewidget.Widget

	corewidget.UsePBuilder(node, newAppState, modules)
	children = []corewidget.Widget{
		sidebar(),
		content(),
	}

	// children = []corewidget.Widget{
	// 	dom.Label(
	// 		dom.AttrStyleParse(
	// 			"width", "30px",
	// 			"height", "30px",
	// 			// "font-style", "italic",
	// 			"font-family", `Helvetica, "Segoe UI", "Apple Color Emoji","Segoe UI Emoji"`,
	// 			"content-align", "center",
	// 			"transform", "scale(8.0)",
	// 			"background-color", "fade-out(green, 60%)",
	// 			"text-decoration", "underline",
	// 		),
	// 		// dom.AttrLabelData("A"),
	// 		dom.AttrLabelData("😀"),
	// 		// dom.AttrLabelData("A😀😳B"),
	// 	),
	// }

	//
	//
	//children = []corewidget.Widget{
	//	dom.Image(
	//		dom.AttrImageSrc(imgsrc.DPI(imgsrc.FilePath(fileutil.ExpandHomeDir("~/Desktop/1.png")), 2)),
	//		dom.AttrStyleParse(
	//			"width", "100%",
	//			"height", "100%",
	//			"object-fit", "cover",
	//			"object-position", "center",
	//		),
	//	),
	//}

	return dom.Style(corebase.PkgFile("app.scss"))(
		dom.Div(dom.Id("app")).Children(
			children...,
		),
	)
})

func Content() corewidget.Widget {
	return _App([]appModule{
		{Title: "Basic", Content: basic.App()},
		{Title: "FileBrowser", Content: filebrowser.App()},
		{Title: "Edit", Content: edit.App()},
		{Title: "Tray", Content: tray.App()},
		{Title: "HotKey", Content: hotkey.App()},
		{Title: "Menu", Content: menu.App()},
		{Title: "Window", Content: window.App()},
		{Title: "Router", Content: router.App()},
		{Title: "I18n", Content: i18n.App()},
		{Title: "Popup", Content: popup.App()},
	})
}
