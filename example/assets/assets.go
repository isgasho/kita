package assets

import (
	"github.com/uiez/uikit/core/corebase"
	"github.com/uiez/uikit/core/coreglobal"
	"github.com/uiez/uikit/modules/locales/i18n"
)

//#uikit asset: i18n fonts images

var AppI18n = i18n.File(corebase.PkgFile("i18n/app.json"))

func init() {
	coreglobal.AddFonts(
		corebase.PkgFile("fonts/NotoColorEmoji.ttf"),
		corebase.PkgFile("fonts/seguiemj.ttf"),
		corebase.PkgFile("fonts/SourceCodePro-Regular.ttf"),
	)
}
