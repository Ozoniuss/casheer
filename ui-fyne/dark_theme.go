package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type sleekDark struct{}

func (s *sleekDark) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 22, G: 24, B: 28, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 38, G: 43, B: 50, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 90, G: 90, B: 90, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 234, G: 236, B: 239, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 48, G: 54, B: 61, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 32, G: 36, B: 42, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 78, G: 156, B: 255, A: 255}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 60, G: 120, B: 220, A: 120}
	case theme.ColorNameShadow:
		return color.NRGBA{0, 0, 0, 150}
	default:
		return theme.DarkTheme().Color(n, v)
	}
}
func (s *sleekDark) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }
func (s *sleekDark) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DarkTheme().Icon(n) }
func (s *sleekDark) Size(n fyne.ThemeSizeName) float32 {
	switch n {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameScrollBar:
		return 6
	default:
		return theme.DarkTheme().Size(n)
	}
}
