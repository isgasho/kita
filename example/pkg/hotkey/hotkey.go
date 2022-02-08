package hotkey

import (
	"github.com/uiez/uikit/app"
	"github.com/uiez/uikit/components/hotkey"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/coreevent"
	"github.com/uiez/uikit/core/corert"
	"github.com/uiez/uikit/core/coreui"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
)

func useHotkeys(node corewidget.ComponentNode, enableCtrlD, enableF1 bool) bool {
	hotkeyGroup := corert.NewOptional(hotkey.UseGroup(node))

	var window *app.Window
	coreui.GetModuleInto(node.View(), &window)
	corewidget.UseEffect(node, func() corebase.CancelFunc {
		if !hotkeyGroup.OK() {
			return func() {}
		}
		var shortcuts []coreevent.Shortcut
		if enableF1 {
			shortcuts = append(shortcuts, coreevent.Shortcut{
				Key:  coreevent.KeyF1,
				Mods: 0,
			})
		}
		if enableCtrlD {
			shortcuts = append(shortcuts, coreevent.Shortcut{
				Key:  coreevent.KeyD,
				Mods: coreevent.ModControl,
			})
		}
		hotkeyGroup.Value().Update(shortcuts, func(action string) {
			window.SetVisible(!window.State().Visible)
		})
		return nil
	}, enableCtrlD, enableF1)
	corewidget.UseEffect(node, func() corebase.CancelFunc {
		if hotkeyGroup.OK() {
			return func() {
				hotkeyGroup.Value().Destroy()
			}
		}
		return nil
	})
	return hotkeyGroup.OK()
}

var App = corewidget.Component("Hotkey", func(node corewidget.ComponentNode) corewidget.Widget {
	disableOrEnable := func(enabled bool) string {
		if !enabled {
			return "enable"
		}
		return "disable"
	}
	enableCtrlD := corewidget.UseValue(node, false)
	enableF1 := corewidget.UseValue(node, false)
	supportHotkey := useHotkeys(node, enableCtrlD.Get(), enableF1.Get())

	errorMsg := corert.Conditional(supportHotkey, "", "hotkey service not available")

	return dom.Style(corebase.PkgFile("hotkey.scss"))(
		dom.Div(dom.Id("pkg")).Compose(func(emit corewidget.ComposeEmitFunc) {
			if errorMsg != "" {
				dom.Label(dom.LabelText(errorMsg)).Emit(emit)
			}
			dom.Div(dom.Id("actions")).Children(
				dom.Button(
					dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
						enableCtrlD.Set(!enableCtrlD.Get())
					}),
				).Child(dom.Label(dom.LabelText(disableOrEnable(enableCtrlD.Get())+" Ctrl+D to show/hide window"))),
				dom.Button(
					dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
						enableF1.Set(!enableF1.Get())
					}),
				).Child(dom.Label(dom.LabelText(disableOrEnable(enableF1.Get())+" F1 to show/hide window"))),
			).Emit(emit)
		}),
	)
})
