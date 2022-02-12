package popup

import (
	"encoding/base64"
	"io"
	"math/rand"
	"time"

	"github.com/uiez/uikit/components/tooltip"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/coregfx"
	"github.com/uiez/uikit/core/coreui"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/dom/styles"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
	"github.com/uiez/uikit/modules/utils/ioutil2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randColor() coregfx.Color {
	v := rand.Uint32()
	var c coregfx.Color
	c.R = uint8(v >> 24 & 0xff)
	c.G = uint8(v >> 16 & 0xff)
	c.B = uint8(v >> 8 & 0xff)
	c.A = 0xff
	return c
}

func randText() string {
	b := make([]byte, 16)

	_, _ = io.ReadFull(ioutil2.ReaderFunc(rand.Read), b)
	return base64.StdEncoding.EncodeToString(b)
}

func randDialog(ele coredom.Element) corewidget.Widget {
	return dom.Div(
		dom.Class("dialog"),
		dom.InlineStyle(styles.PropBackgroundColor.NewRuleProp(randColor())),
		dom.LogicalParent(ele),
	).Children(
		dom.Label(
			dom.LabelText(randText()),
		),
	)
}

func randTip(ele coredom.Element, class string, text string) corewidget.Widget {
	return dom.Label(dom.Class(class), dom.LabelText(text), dom.LogicalParent(ele))
}

type buttonItem struct {
	OnTap func(node corewidget.Node)
	Text  string
}

func buttonList() []buttonItem {
	return []buttonItem{
		{
			func(node corewidget.Node) {
				dom.PopupPush(node,
					[]corewidget.Widget{randDialog(node.Element())},
					dom.PopupOptModal(0.3),
					dom.PopupOptDismissOnEscapeOrBack(),
					dom.PopupOptDismissOnPointerDownOutside(),
				)
			},
			"ShowDialog(dismiss by ESC/click background)",
		},
		{
			func(node corewidget.Node) {
				handle := dom.PopupPush(node, []corewidget.Widget{randDialog(node.Element())},
					dom.PopupOptModal(0.3),
					dom.PopupOptDismissOnEscapeOrBack(),
					dom.PopupOptDismissOnPointerDownOutside(),
				)
				handle.OnClose(func(kind dom.PopupCloseKind) bool {
					if kind.ByExternalEvent() {
						tipHandle := dom.PopupPush(node, []corewidget.Widget{
							randTip(node.Element(), "notification", "Dialog dismissed"),
						})
						coreui.AutoCancel(node.View(), time.Second*2, tipHandle.Dismiss)
					}
					return true
				})
			},
			"ShowDialog(dismiss by ESC/click background, show tips)",
		},
		{
			func(node corewidget.Node) {
				dialogHandle := dom.PopupPush(node,
					[]corewidget.Widget{randDialog(node.Element())},
					dom.PopupOptModal(0.3),
				)
				coreui.AutoCancel(node.View(), time.Second*2, dialogHandle.Dismiss)
			},
			"ShowDialog(auto dismiss in 2 sec)",
		},
		{
			func(node corewidget.Node) {
				dialogHandle := dom.PopupPush(node, []corewidget.Widget{randDialog(node.Element())},
					dom.PopupOptModal(0),
					dom.PopupOptDismissOnEscapeOrBack(),
					dom.PopupOptDismissOnPointerDownOutside(),
				)
				var (
					times          int
					againTimer     corebase.CancelFunc
					againTipHandle dom.OptionalPopupHandle
				)
				dialogHandle.OnClose(func(kind dom.PopupCloseKind) bool {
					times++
					if times >= 2 || !kind.Cancellable() {
						againTimer.SafeCancel()
						againTipHandle.Dismiss()
						return true
					}

					againTipHandle.Reset(dom.PopupPush(node, []corewidget.Widget{randTip(node.Element(), "notification", "Again to dismiss dialog in 2 sec")}))
					againTimer.SafeCancel()
					const againTimeout = time.Second * 2
					againTimer = node.View().AddTimer(againTimeout, false, func(_ time.Time) {
						times = 0
						againTimer = nil
						againTipHandle.Dismiss()
					})
					return false
				})
			},
			"ShowDialog(dismiss by twice ESC/click background)",
		},
		{
			func(node corewidget.Node) {
				mainHandle := dom.PopupPush(node, []corewidget.Widget{randDialog(node.Element())}, dom.PopupOptModal(0),
					dom.PopupOptDismissOnEscapeOrBack(),
					dom.PopupOptDismissOnPointerDownOutside())

				var confirmHandle dom.OptionalPopupHandle
				confirm := dom.Div(dom.Class("confirm-panel")).Children(
					randDialog(node.Element()),
					dom.Button(
						dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
							mainHandle.Dismiss()
							confirmHandle.Dismiss()
						}),
					).Child(dom.Label(dom.LabelText("OK"))),
				)
				mainHandle.OnClose(func(kind dom.PopupCloseKind) bool {
					if !kind.Cancellable() {
						return true
					}
					confirmHandle.Reset(dom.PopupPush(node, []corewidget.Widget{confirm},
						dom.PopupOptModal(0),
						dom.PopupOptDismissOnEscapeOrBack(),
						dom.PopupOptDismissOnPointerDownOutside(),
					))
					return false
				})
			},
			"ShowDialog(dismiss by another confirm dialog)",
		},
	}
}

var tooltipButton = corewidget.PComponent("TooltipButton", func(node corewidget.ComponentNode, item buttonItem) corewidget.Widget {
	return tooltip.Text("Tooltip").Child(
		dom.Button(
			dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
				item.OnTap(node)
			}),
		).Child(dom.Label(dom.LabelText(item.Text))),
	)
})

func buildItems(buttons []buttonItem) []corewidget.Widget {
	items := make([]corewidget.Widget, 0, len(buttons))
	for i := range buttons {
		items = append(items, tooltipButton(buttons[i]))
	}
	return items
}

var App = corewidget.Component("Popup", func(node corewidget.ComponentNode) corewidget.Widget {
	buttons := buttonList()
	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("popup.scss"))
	return dom.Div(dom.Stylesheets(stylesheets), dom.Id("pkg")).Children(
		dom.PopupContainer(dom.Div(dom.Class("panel")).Children(buildItems(buttons)...)),
		dom.Spacer(),
		dom.PopupContainer(dom.Div(dom.Class("panel")).Children(buildItems(buttons)...)),
	)
})
