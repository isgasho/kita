package i18n

import (
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/example/assets"
	"github.com/uiez/uikit/modules/locales/i18n"
)

var App = corewidget.Component("I18n", func(node corewidget.ComponentNode) corewidget.Widget {
	text := i18n.UseTranslate(node, assets.AppI18n.Tr("hello"))

	return dom.Span(dom.SpanText(text))
})
