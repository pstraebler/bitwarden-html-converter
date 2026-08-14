package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Bitwarden HTML Converter")
	myWindow.Resize(fyne.NewSize(600, 400))

	var inputPath string
	var outputPath string

	inputLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("No destination selected")
	statusLabel := widget.NewLabel("")

	inputButton := widget.NewButton("Select JSON File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			inputPath = reader.URI().Path()
			inputLabel.SetText(filepath.Base(inputPath))
			reader.Close()
		}, myWindow)
	})

	outputButton := widget.NewButton("HTML Destination", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if writer == nil {
				return
			}
			outputPath = writer.URI().Path()
			if filepath.Ext(outputPath) != ".html" {
				outputPath += ".html"
			}
			outputLabel.SetText(filepath.Base(outputPath))
			writer.Close()
		}, myWindow)
	})

	convertButton := widget.NewButton("Convert", func() {
		if inputPath == "" {
			dialog.ShowError(fmt.Errorf("Please select a JSON file"), myWindow)
			return
		}
		if outputPath == "" {
			dialog.ShowError(fmt.Errorf("Please select a destination"), myWindow)
			return
		}

		statusLabel.SetText("Converting...")
		err := ConvertBitwardenToHTML(inputPath, outputPath)
		if err != nil {
			statusLabel.SetText("Conversion error")
			dialog.ShowError(err, myWindow)
			return
		}
		statusLabel.SetText("Conversion successful!")
		dialog.ShowInformation("Success", "HTML file created successfully!", myWindow)
	})

	quitButton := widget.NewButton("Quit", func() {
		myApp.Quit()
	})

	content := container.NewVBox(
		widget.NewLabel("Bitwarden JSON to HTML Converter"),
		widget.NewSeparator(),
		inputButton,
		inputLabel,
		widget.NewSeparator(),
		outputButton,
		outputLabel,
		widget.NewSeparator(),
		convertButton,
		statusLabel,
		widget.NewSeparator(),
		quitButton,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
