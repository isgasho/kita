package edit

import (
	"fmt"
	"io/ioutil"

	"github.com/uiez/uikit/components/dnd"
	"github.com/uiez/uikit/components/filechoose"
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coredom"
	"github.com/uiez/uikit/core/corewidget"
	"github.com/uiez/uikit/dom"
	"github.com/uiez/uikit/modules/events/gestures/recognizers"
	"github.com/uiez/uikit/modules/locales/i18n"
	"github.com/uiez/uikit/modules/services"
	"github.com/uiez/uikit/modules/texts"
	"github.com/uiez/uikit/modules/utils/fileutil"
	"github.com/uiez/uikit/modules/utils/mimeutil"
	"github.com/uiez/uikit/modules/utils/unsafeptr"
)

func init() {
	mimeutil.AddExtensionType(".go", "text/go")
	mimeutil.AddExtensionType(".md", "text/markdown")
}

var App = corewidget.Component("Edit",
	func(node corewidget.ComponentNode) corewidget.Widget {
		message := corewidget.UseValue(node, "")
		text := corewidget.UseValue[texts.Text](node, texts.Runes(""))
		loadFile := func(path string) {
			message.Set("loading file...")
			go func() {
				data, err := ioutil.ReadFile(path)

				node.View().AddTask(func() {
					if err == nil {
						text.Set(texts.Runes(unsafeptr.BytesToString(data)))
						message.Set("file loaded:" + path)
					} else {
						message.Set(fmt.Sprint("read file failed:", path, err))
					}
				})
			}()
		}

		saveFile := func(path string) {
			data := texts.CopyTextDataAsString(text.Get())
			message.Set("saving file...")
			go func() {
				err := ioutil.WriteFile(path, unsafeptr.StringToBytes(data), 0o644)

				node.View().AddTask(func() {
					if err != nil {
						message.Set(fmt.Sprint("write file failed:", path, err))
					} else {
						message.Set("file saved:" + path)
					}
				})
			}()
		}

		exts := []services.FileFilter{
			{Description: "Text", MIMEs: []string{"text/plain"}},
			{Description: "Markdown", MIMEs: []string{"text/markdown"}},
			{Description: "Go", MIMEs: []string{"text/go"}},
		}
		fileChooser := filechoose.UseState(node)
		onFileOpen := func(ele coredom.Element, details recognizers.TapDetails) {
			fileChooser.Open(i18n.Str("Choose file to open"),
				exts,
				services.FileDialogAllowFile|services.FileDialogAllowMulti,
				func(selected bool, paths []string) {
					if selected && len(paths) > 0 {
						loadFile(paths[0])
					}
				})
		}
		onFileSave := func(ele coredom.Element, details recognizers.TapDetails) {
			fileChooser.Save(i18n.Str("Choose file to save"), "data.go", exts, func(selected bool, path string) {
				if selected {
					saveFile(path)
				}
			})
		}
		_ = onFileOpen
		_ = onFileSave
		dragDropHandler := dnd.UseDNDHandler(node, func(paths []string) bool {
			return len(paths) == 1 && fileutil.IsFile(paths[0])
		}, func(paths []string) bool {
			loadFile(paths[0])
			return true
		})

		return dom.Style(corebase.PkgFile("edit.scss"))(
			dom.Div(dom.Id("pkg")).Children(
				dom.Div(dom.Id("dialog-actions")).Children(
					dom.Button(
						dom.Class("dialog-action"),
						dom.OnTap(onFileOpen),
					).Child(dom.Label(dom.LabelText("Open"))),

					dom.Button(
						dom.Class("dialog-action"),
						dom.OnTap(onFileSave),
					).Child(dom.Label(dom.LabelText("Save"))),

					dom.Label(dom.LabelText(message.Get())),
				),

				dom.TextArea(
					dom.OnDragDrop(dragDropHandler),
					dom.TextAreaPlaceholder("File Content(support drag drop file here)"),
					dom.TextAreaValue(text.Get()),
					dom.TextAreaOnChange(text.Set),
				),
			),
		)
	},
)
