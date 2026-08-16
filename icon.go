package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed icon.png
var iconData []byte

func getAppIcon() fyne.Resource {
	return fyne.NewStaticResource("icon.png", iconData)
}
