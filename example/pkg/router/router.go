package router

import (
	"encoding/base64"
	"io"
	"math/rand"
	"time"

	"github.com/uiez/uikit/components/router"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/coregfx"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
	"github.com/uiez/uikit/modules/gfx/canvaskit"
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
	b := make([]byte, 32)
	_, _ = io.ReadFull(ioutil2.ReaderFunc(rand.Read), b)
	return base64.StdEncoding.EncodeToString(b)
}

func App() corewidget.Widget {
	routeColor := "/color"
	routeText := "/text"
	actions := func(state router.State, route router.RouteBuildInfo) []corewidget.Widget {
		buttons := []struct {
			OnTap   func()
			Text    string
			Disable bool
		}{
			{
				func() { state.Push(routeColor, randColor(), nil) },
				"PushColor",
				route.Name == routeColor || route.Name == "",
			},
			{
				func() { state.Push(routeText, "Hello "+randText(), nil) },
				"PushText",
				route.Name == routeText,
			},
			{
				func() { state.Pop(nil) },
				"Pop",
				state.Count() <= 1,
			},
		}

		items := make([]corewidget.Widget, len(buttons))
		for i := range items {
			b := buttons[i]
			items[i] = dom.Button(
				dom.Disabled(b.Disable),
				dom.Class(b.Text),
				dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
					b.OnTap()
				}),
			).Child(
				dom.Label(dom.LabelText(b.Text)),
			)
		}
		return items
	}

	routes := router.NamedRouteBuilder{
		DefaultRouteName: routeColor,
		Routes: map[string]func(params any) corewidget.Widget{
			routeColor: func(params any) corewidget.Widget {
				col, ok := params.(coregfx.Color)
				if !ok {
					col = randColor()
				}
				return dom.Canvas(dom.CanvasPainter(canvaskit.ColorPainterBuilder(col)))
			},
			routeText: func(params any) corewidget.Widget {
				text, _ := params.(string)
				if text == "" {
					text = "Default"
				}
				return dom.Label(dom.LabelText(text))
			},
		},
	}
	return router.Router(
		routes.Build,
		func(state router.State, info router.RouteBuildInfo, routerView corewidget.Widget) corewidget.Widget {
			return dom.Style(corebase.PkgFile("router.scss"))(
				dom.Div(dom.Id("pkg")).Children(
					dom.Div(dom.Id("router-view")).Children(routerView),
					dom.Div(dom.Id("router-actions")).Children(
						actions(state, info)...,
					),
				),
			)
		},
	)
}
