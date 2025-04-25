package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"securezip_gui/zipcrypto"
	"strings"
)

func RunApp() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SecureZip")

	title := widget.NewLabelWithStyle("SecureZip", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	inputLabel := widget.NewLabel("Choisir un fichier :")
	inputPath := widget.NewEntry()
	browseBtn := widget.NewButtonWithIcon("Parcourir", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				inputPath.SetText(uri.URI().Path())
			}
		}, myWindow)
	})

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Mot de passe...")

	status := widget.NewLabel("")
	status.Wrapping = fyne.TextWrapWord

	encryptBtn := widget.NewButtonWithIcon("Chiffrer", theme.ConfirmIcon(), func() {
		folder := inputPath.Text
		pass := password.Text
		if len(pass) < 6 {
			dialog.ShowError(fmt.Errorf("Mot de passe trop court"), myWindow)
			return
		}

		outFile := folder + ".enc"
		err := zipcrypto.EncryptFile(folder, outFile, pass)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		status.SetText("Chiffré avec succès : " + outFile)
	})

	decryptBtn := widget.NewButtonWithIcon("Déchiffrer", theme.ContentUndoIcon(), func() {
		file := inputPath.Text
		pass := password.Text
		if len(pass) < 6 {
			dialog.ShowError(fmt.Errorf("Mot de passe trop court"), myWindow)
			return
		}

		if !strings.HasSuffix(file, ".enc") {
			dialog.ShowError(fmt.Errorf("Le fichier doit avoir l'extension .enc"), myWindow)
			return
		}

		status.SetText("Déchiffrement en cours...")

		outFile := file[:len(file)-4]
		err := zipcrypto.DecryptFile(file, outFile, pass)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Erreur lors du déchiffrement : %v", err), myWindow)
			status.SetText("Erreur lors du déchiffrement.")
			return
		}
		status.SetText("Déchiffré avec succès : " + outFile)
	})

	content := container.NewVBox(
		title,
		inputLabel,
		container.NewBorder(nil, nil, nil, browseBtn, inputPath),
		widget.NewLabel("Mot de passe :"),
		password,
		container.NewHBox(encryptBtn, decryptBtn),
		status,
	)

	myWindow.SetContent(container.NewPadded(content))
	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.ShowAndRun()
}