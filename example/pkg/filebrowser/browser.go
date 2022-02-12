package filebrowser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/uiez/uikit/components/toast"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/dom/styles"
	"github.com/uiez/uikit/modules/css/csstypes"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
	"github.com/uiez/uikit/modules/utils/mimeutil"
)

type treeNode struct {
	name string // display name

	path  string // relative path to tree dir
	isDir bool

	open     bool
	children []*treeNode
}

type browserState struct {
	corewidget.StateBase
	node corewidget.Node
	dir  string
	root *treeNode

	openedFile        string
	openedFileContent string
}

func newBrowserState(node corewidget.ComponentNode) *browserState {
	homeDir, _ := os.UserHomeDir()
	dir := homeDir
	s := &browserState{
		dir:  dir,
		node: node,
		root: &treeNode{
			name:  dir,
			path:  ".",
			isDir: true,
		},
	}

	s.toggleDir(s.root)
	return s
}

func (s *browserState) toggleDir(node *treeNode) {
	if !node.isDir {
		return
	}

	if node.open {
		node.open = false
		node.children = nil
		if node.path == "." || strings.HasPrefix(s.openedFile, node.path+"/") {
			s.openedFile = ""
			s.openedFileContent = ""
		}
	} else {
		node.open = true
		items, _ := ioutil.ReadDir(filepath.Join(s.dir, node.path))
		sort.Slice(items, func(i, j int) bool {
			if items[i].IsDir() != items[j].IsDir() {
				return items[i].IsDir()
			}
			return items[i].Name() < items[j].Name()
		})
		node.children = make([]*treeNode, 0, len(items))
		for _, v := range items {
			name := v.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			node.children = append(node.children, &treeNode{
				name:  v.Name(),
				path:  filepath.Join(node.path, v.Name()),
				isDir: v.IsDir(),
			})
		}
	}
	s.Notify()
}

func (s *browserState) openFile(node *treeNode) {
	path := filepath.Join(s.dir, node.path)
	mime := mimeutil.TypeByExtension(filepath.Ext(path))
	if !strings.HasPrefix(mime, "text/") {
		toast.Show(s.node, "unsupported file type")
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		toast.Show(s.node, "open file failed:"+err.Error())
		return
	}
	s.openedFile = node.path
	s.openedFileContent = string(data)

	s.Notify()
}

func buildTreeList(state *browserState, widgets []corewidget.Widget, node *treeNode, level int) []corewidget.Widget {
	attrs := []coredom.Attribute{
		dom.Class("item", dom.OptionClass(node.path == state.openedFile, "is-selected")),
		dom.InlineStyle(styles.PropPaddingLeft.NewComputableRuleProp(csstypes.NewEm(float64(level)))),
	}
	if !node.isDir {
		attrs = append(attrs, dom.OnDoubleTap(func(ele coredom.Element) {
			state.openFile(node)
		}))
	}
	widgets = append(widgets,
		dom.Div.Key(node.path)(
			attrs...,
		).Compose(func(emit corewidget.ComposeEmitFunc) {
			if node.isDir {
				dom.Icon(
					dom.OnTap(func(ele coredom.Element, details recognizers.TapDetails) {
						state.toggleDir(node)
					}),
					dom.Class("state", dom.OptionClass(node.open, "is-open")),
					dom.IconCodepoint('>'),
				).Emit(emit)
			}

			dom.Label(
				dom.Class(dom.OptionClass(node.isDir, "is-dir")),
				dom.LabelText(node.name),
			).Emit(emit)
		}),
	)
	if node.isDir && node.open {
		for _, c := range node.children {
			widgets = buildTreeList(state, widgets, c, level+1)
		}
	}
	return widgets
}

var appSidebar = corewidget.Component("Sidebar",
	func(node corewidget.ComponentNode) corewidget.Widget {
		bs := corewidget.UseBuilderConsumer(node, newBrowserState).Value()
		return dom.Div(dom.Class("filetree")).Children(buildTreeList(bs, nil, bs.root, 0)...)
	},
)

var appContent = corewidget.Component("Content",
	func(node corewidget.ComponentNode) corewidget.Widget {
		bs := corewidget.UseBuilderConsumer(node, newBrowserState).Value()
		var text string
		if bs.openedFile != "" {
			text = bs.openedFileContent
		} else {
			text = "No file selected."
		}
		return dom.Span(
			dom.Class("file-content",
				dom.OptionClass(bs.openedFile != "", "file-selected"),
			), dom.SpanText(text))
	},
)

var App = corewidget.Component("FileBrowser",
	func(node corewidget.ComponentNode) corewidget.Widget {
		corewidget.UseBuilder(node, newBrowserState)

		stylesheets := dom.UseStylesheet(node, corebase.PkgFile("browser.scss"))
		return dom.Div(dom.Stylesheets(stylesheets), dom.Id("pkg")).Children(
			appSidebar(),
			appContent(),
		)
	},
)
