package menu

import (
	"github.com/uiez/uikit/components/menu"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coreevent"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/modules/services"
)

func keyMenu() services.MenuItems {
	n := coreevent.MaxKey - coreevent.MinKey + 1
	items := make(services.MenuItems, 0, n)
	for i := coreevent.Key(0); i < n; i++ {
		k := i + coreevent.MinKey
		items = append(items, services.MenuItem{
			Label:         k.String(),
			Action:        k.String(),
			KeyEquivalent: k,
		})
	}
	return items
}

func keyMenuModifier() services.MenuItems {
	mod := coreevent.ModControl
	prefix := "Control+"
	menu := keyMenu()
	for i := range menu {
		menu[i].Label = prefix + menu[i].Label
		menu[i].KeyEquivalentModifiers = mod
		menu[i].Action = prefix + menu[i].Action
	}
	return menu
}

func buildMenubar(checked bool) services.MenuItems {
	return services.MenuItems{
		{
			Label: "App",
			Submenu: services.MenuItems{
				{
					Label: "Group",
					Submenu: services.MenuItems{
						{
							Type:                   services.MenuItemTypeChecked,
							Label:                  "Toggle",
							Action:                 "toggle",
							Checked:                checked,
							KeyEquivalent:          coreevent.KeyF,
							KeyEquivalentModifiers: coreevent.ModCommandOrControl | coreevent.ModOptionOrAlt,
						},
					},
				},
			},
		},
		{
			Label:   "Key",
			Submenu: keyMenu(),
		},
		{
			Label:   "KeyModifiers",
			Submenu: keyMenuModifier(),
		},
	}
}

func buildContextMenu() services.MenuItems {
	return keyMenu()
}

var App = corewidget.Component("Menu", func(node corewidget.ComponentNode) corewidget.Widget {
	checked := corewidget.UseValue(node, false)
	action := corewidget.UseValue(node, "")
	actionSource := corewidget.UseValue(node, "")

	text := corewidget.UseMemo(node, func() string {
		var text string
		if action.Get() == "" {
			text = "action: no action"
		} else {
			text = "action:" + actionSource.Get() + "," + action.Get()
		}
		return text
	}, action.Get(), actionSource.Get())

	onAction := func(source, act string) {
		actionSource.Set(source)
		action.Set(act)
		switch act {
		case "toggle":
			checked.Set(!checked.Get())
		}
	}

	menu.UseWindowMenubar(node, buildMenubar(checked.Get()), func(action string) {
		onAction("menubar", action)
	})

	stylesheets := dom.UseStylesheet(node, corebase.PkgFile("menu.scss"))
	return dom.Div(dom.Stylesheets(stylesheets)).Children(
		dom.Label(dom.LabelText(text)),
		dom.Button(
			menu.AttrSecondaryTapContextMenu(buildContextMenu, func(action string) {
				onAction("contextMenu", action)
			}),
		).Child(dom.Label(dom.LabelText("show context menu"))),
	)
})
