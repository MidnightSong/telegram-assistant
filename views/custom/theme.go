package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type Theme struct{}

// Color 覆盖禁用状态下的颜色设置
func (Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameDisabled {
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255} //
	}
	return theme.DefaultTheme().Color(name, variant)
}

// Icon 返回图标（未自定义，使用默认主题）
func (Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Font 返回字体（未自定义，使用默认主题）
func (Theme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Size 返回字体大小（未自定义，使用默认主题）
func (Theme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
