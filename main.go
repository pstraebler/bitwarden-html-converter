package main

import (
	"fmt"
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const version = "1.1"

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Bitwarden HTML Converter")
	myWindow.Resize(fyne.NewSize(700, 600))

	var inputPath string
	var outputPath string

	inputLabel := widget.NewLabel("No file selected")
	outputLabel := widget.NewLabel("No destination selected")
	statusLabel := widget.NewLabel("")

	// Field selection checkboxes with defaults
	checkType := widget.NewCheck("Type", nil)
	checkType.Checked = true

	checkName := widget.NewCheck("Name", nil)
	checkName.Checked = true

	checkUsername := widget.NewCheck("Username/Login", nil)
	checkUsername.Checked = true

	checkPassword := widget.NewCheck("Password", nil)
	checkPassword.Checked = true

	checkNotes := widget.NewCheck("Notes", nil)
	checkNotes.Checked = true

	checkURL := widget.NewCheck("URL", nil)
	checkURL.Checked = false

	checkFavorite := widget.NewCheck("Show Favorite Star", nil)
	checkFavorite.Checked = false

	checkTOTP := widget.NewCheck("TOTP", nil)
	checkTOTP.Checked = false

	checkFolder := widget.NewCheck("Folder", nil)
	checkFolder.Checked = false

	checkOrganization := widget.NewCheck("Organization", nil)
	checkOrganization.Checked = false

	checkCreationDate := widget.NewCheck("Creation Date", nil)
	checkCreationDate.Checked = false

	checkModificationDate := widget.NewCheck("Modification Date", nil)
	checkModificationDate.Checked = false

	checkPasswordRevision := widget.NewCheck("Password Revision Date", nil)
	checkPasswordRevision.Checked = false

	checkCustomFields := widget.NewCheck("Custom Fields", nil)
	checkCustomFields.Checked = false

	// Grouping options
	groupingLabel := widget.NewLabel("Group entries by:")
	groupingSelect := widget.NewSelect([]string{"None", "Type", "Folder"}, nil)
	groupingSelect.SetSelected("None")

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

		// Build ExportFields struct from checkboxes
		fields := ExportFields{
			Type:              checkType.Checked,
			Name:              checkName.Checked,
			Username:          checkUsername.Checked,
			Password:          checkPassword.Checked,
			Notes:             checkNotes.Checked,
			URL:               checkURL.Checked,
			Favorite:          checkFavorite.Checked,
			TOTP:              checkTOTP.Checked,
			Folder:            checkFolder.Checked,
			Organization:      checkOrganization.Checked,
			CreationDate:      checkCreationDate.Checked,
			ModificationDate:  checkModificationDate.Checked,
			PasswordRevision:  checkPasswordRevision.Checked,
			CustomFields:      checkCustomFields.Checked,
		}

		// Determine grouping mode
		var grouping GroupingMode
		switch groupingSelect.Selected {
		case "Type":
			grouping = GroupByType
		case "Folder":
			grouping = GroupByFolder
		default:
			grouping = GroupNone
		}

		options := ExportOptions{
			Fields:   fields,
			Grouping: grouping,
		}

		statusLabel.SetText("Converting...")
		err := ConvertBitwardenToHTML(inputPath, outputPath, options)
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

	// Group checkboxes into columns for better layout
	leftColumn := container.NewVBox(
		checkType,
		checkName,
		checkUsername,
		checkPassword,
		checkNotes,
		checkURL,
		checkFavorite,
	)

	rightColumn := container.NewVBox(
		checkTOTP,
		checkFolder,
		checkOrganization,
		checkCreationDate,
		checkModificationDate,
		checkPasswordRevision,
		checkCustomFields,
	)

	fieldsContainer := container.NewHBox(leftColumn, rightColumn)

	// Version label with discrete styling
	versionText := canvas.NewText(fmt.Sprintf("Bitwarden HTML Converter v%s", version), color.RGBA{R: 128, G: 128, B: 128, A: 255})
	versionText.TextSize = 10
	versionText.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		widget.NewLabel("Bitwarden JSON to HTML Converter"),
		widget.NewSeparator(),
		inputButton,
		inputLabel,
		widget.NewSeparator(),
		outputButton,
		outputLabel,
		widget.NewSeparator(),
		widget.NewLabel("Select fields to export:"),
		fieldsContainer,
		widget.NewSeparator(),
		groupingLabel,
		groupingSelect,
		widget.NewSeparator(),
		convertButton,
		statusLabel,
		widget.NewSeparator(),
		quitButton,
		widget.NewSeparator(),
		container.NewCenter(versionText),
	)

	myWindow.SetContent(content)

	// Add keyboard shortcuts to quit the application
	ctrlQ := &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}
	altF4 := &desktop.CustomShortcut{KeyName: fyne.KeyF4, Modifier: fyne.KeyModifierAlt}

	myWindow.Canvas().AddShortcut(ctrlQ, func(shortcut fyne.Shortcut) {
		myApp.Quit()
	})
	myWindow.Canvas().AddShortcut(altF4, func(shortcut fyne.Shortcut) {
		myApp.Quit()
	})

	myWindow.ShowAndRun()
}
