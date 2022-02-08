package basic

import (
	"github.com/uiez/uikit/components/colorpicker"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coregfx"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/modules/gfx/colors"
)

type basicItem struct {
	Name string

	Content corewidget.Widget
}

func buildCheckbox(node corewidget.ComponentNode, checkboxState *corewidget.Ref[corebase.SortedStrings], disabled bool, vals ...string) corewidget.Widget {
	hasAll := true
	hasOne := false
	for _, v := range vals {
		if !checkboxState.Current.Has(v) {
			hasAll = false
		} else {
			hasOne = true
		}
	}
	return dom.Checkbox(
		dom.CheckboxChecked(hasAll),
		dom.CheckboxIndeterminate(hasOne && !hasAll),
		dom.Disabled(disabled),
		dom.CheckboxOnChange(func(checked bool) {
			for _, val := range vals {
				if checked {
					checkboxState.Current.Add(val)
				} else {
					checkboxState.Current.Remove(val)
				}
			}
			node.MarkNeedsBuild()
		}),
	)
}

func buildItems(node corewidget.ComponentNode,
	checkboxState *corewidget.Ref[corebase.SortedStrings],
	radioSeq corewidget.Value[int],
	selectValue corewidget.Value[string]) []basicItem {
	return []basicItem{
		{
			"URL",
			dom.Div(dom.Class("url")).Children(
				dom.Anchor(
					dom.AnchorText("Google"),
					dom.AnchorHref("https://google.com"),
				),
				dom.Anchor(
					dom.Disabled(true),
					dom.AnchorText("Apple"),
					dom.AnchorHref("https://apple.com"),
				),
			),
		},
		{
			"Checkbox",
			dom.Div(dom.Class("checkbox")).Children(
				buildCheckbox(node, checkboxState, false, "A"),
				dom.Label(dom.LabelText("AAAA"), dom.LabelForRelativeBrother(-1)),
				buildCheckbox(node, checkboxState, false, "B"),
				dom.Label(dom.LabelText("BBBB"), dom.LabelForNthBrother(3)),
				buildCheckbox(node, checkboxState, false, "A", "B"),
				dom.Label(dom.LabelText("A+B"), dom.LabelForRelativeBrother(-1)),
				buildCheckbox(node, checkboxState, true, "C"),
				dom.Label(dom.LabelText("CCCC"), dom.LabelForRelativeBrother(-1)),
			),
		},
		{
			"Radio",
			dom.Div(dom.Class("radio")).Children(
				dom.Radio(
					dom.RadioChecked(radioSeq.Get() == 1),
					dom.RadioOnChange(func() { radioSeq.Set(1) }),
				),
				dom.Label(dom.LabelText("AAAA"), dom.LabelForRelativeBrother(-1)),
				dom.Radio(
					dom.RadioChecked(radioSeq.Get() == 2),
					dom.Disabled(true),
					dom.RadioOnChange(func() { radioSeq.Set(2) }),
				),
				dom.Label(dom.LabelText("BBBB"), dom.LabelForNthBrother(3)),
				dom.Radio(
					dom.Class("ring"),
					dom.RadioChecked(radioSeq.Get() == 3),
					dom.RadioOnChange(func() { radioSeq.Set(3) }),
				),
				dom.Label(dom.LabelText("CCCC"), dom.LabelForRelativeBrother(-1)),
			),
		},
		{
			"Select",
			dom.Div(dom.Class("select")).Children(
				dom.Label(dom.LabelText("SelectLabel:"), dom.LabelForRelativeBrother(1)),
				dom.Select(
					dom.SelectPlaceholder("Select a item"),
					dom.SelectValue(selectValue.Get()),
					dom.SelectOnChange(selectValue.Set),
					dom.SelectOptions([]dom.SelectOption{
						{Value: "1", Label: "Item 1", Disabled: false},
						{Value: "2", Label: "Item 22", Disabled: false},
						{Value: "3", Label: "Item 333", Disabled: true},
						{Value: "4", Label: "Item 4444", Disabled: false},
					}...),
				),
			),
		},
		{
			"Input",
			dom.Div(dom.Class("input")).Children(
				dom.Label(dom.LabelText("InputLabel:"), dom.LabelForRelativeBrother(1)),
				dom.Input(
					dom.InputPlaceholder("input"),
				),
			),
		},
		{
			"ColorPicker",
			dom.Div(dom.Class("colorpicker")).Children(
				colorpicker.Input(colors.Green, func(col coregfx.Color) {}),
			),
		},
	}
}

var App = corewidget.Component("Basic",
	func(node corewidget.ComponentNode) corewidget.Widget {
		checkboxState := corewidget.UseRef(node, corebase.SortedStrings{})
		radioSeq := corewidget.UseValue(node, 0)
		selectValue := corewidget.UseValue(node, "")
		items := buildItems(node, checkboxState, radioSeq, selectValue)
		return dom.Style(corebase.PkgFile("basic.scss"))(
			dom.Div(dom.Id("pkg")).Compose(func(emit corewidget.ComposeEmitFunc) {
				dom.Div(dom.Id("nav")).Compose(func(emit corewidget.ComposeEmitFunc) {
					for _, item := range items {
						dom.Anchor(dom.AnchorText(item.Name), dom.AnchorHref("#"+item.Name)).Emit(emit)
					}
				}).Emit(emit)

				dom.Div(dom.Id("content")).Compose(func(emit corewidget.ComposeEmitFunc) {
					for _, item := range items {
						dom.Div(dom.Id(item.Name), dom.Class("item")).Children(
							dom.Label(dom.LabelText(item.Name)),
							item.Content,
						).Emit(emit)
					}
				}).Emit(emit)
			}),
		)
	},
)
