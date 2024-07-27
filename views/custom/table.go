package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type WrapLabel struct {
	widget.BaseWidget
	text string
}

func NewWrapLabel(text string) *WrapLabel {
	w := &WrapLabel{text: text}
	w.ExtendBaseWidget(w)
	return w
}

func (w *WrapLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(w.text, theme.ForegroundColor())
	text.TextStyle = fyne.TextStyle{Italic: true}
	return &wrapLabelRenderer{text}
}

func (w *WrapLabel) SetText(text string) {
	w.text = text
}

type wrapLabelRenderer struct {
	text *canvas.Text
}

func (r *wrapLabelRenderer) Layout(size fyne.Size) {
	r.text.TextSize = 14
	r.text.Resize(size)
}

func (r *wrapLabelRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *wrapLabelRenderer) Refresh() {
	canvas.Refresh(r.text)
}

func (r *wrapLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}

func (r *wrapLabelRenderer) Destroy() {}

/*themes := container.NewGridWithColumns(2,
widget.NewButton("Dark", func() {
	a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
}),
widget.NewButton("Light", func() {
	a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
}),
)
*/
