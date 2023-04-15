package gui

import (
	"fmt"

	"hqdragondownloader/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func App() {

	data := utils.HQs{
		Links: []string{},
		Names: []string{},
	}

	dataCaps := utils.HQCaps{
		Links: []string{},
		Caps:  []string{},
	}

	nameHQ := ""

	path2Output := "."

	linkHQ := ""

	linkHQCap := ""

	HQCap := ""

	myApp := app.New()
	myWindow := myApp.NewWindow("HQDragonDownloader")

	myWindow.Resize(fyne.NewSize(800, 600))

	selectInputCap := widget.NewSelect(dataCaps.Caps, func(s string) {
		idxCap := -1
		for i, item := range dataCaps.Caps {
			if item == s {
				idxCap = i
				break
			}
		}

		HQCap = s
		linkHQCap = dataCaps.Links[idxCap]
	})

	selectInputCap.PlaceHolder = "Capítulo da HQ"

	label := widget.NewLabel("Pesquise por uma HQ!")

	labelResult := widget.NewLabel("")

	labelResult.Alignment = fyne.TextAlignCenter

	label.Alignment = fyne.TextAlignCenter

	w := fyne.CurrentApp().NewWindow("HQPageToDownload")

	selectPath := widget.NewButton("Salvar onde?", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			//children, err := list.List()

			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			//out := fmt.Sprintf("Folder %s (%d children):\n%s", list.Name(), len(children), list.String())
			dialog.ShowInformation("Pasta selecionada:", list.Path(), w)

			path2Output = list.Path()

		}, w)
	})

	newWindow := container.NewGridWithRows(4, selectInputCap, widget.NewButton("Baixar", func() {
		// labelResult.Text = "Baixando..."
		// labelResult.Refresh()

		if selectInputCap.Selected == "" {
			labelResult.Text = "Selecione um capítulo."
			labelResult.Refresh()
			return
		}

		if HQCap == "All" {
			utils.DownloadHQ(dataCaps.Links, nameHQ, labelResult, path2Output)
		} else {
			utils.DownloadHQ([]string{linkHQCap}, nameHQ, labelResult, path2Output)
		}

		labelResult.Text = "Pronto"
		labelResult.Refresh()
	}), labelResult, selectPath)
	//icon := widget.NewIcon(nil)
	button := widget.NewButtonWithIcon("", nil, func() {
		w = fyne.CurrentApp().NewWindow("HQPageToDownload")
		dataCaps = utils.GetCaps(linkHQ)
		selectInputCap.Options = dataCaps.Caps
		w.Resize(fyne.NewSize(480, 380))
		w.SetTitle(fmt.Sprintf("HQPageToDownload - %v", nameHQ))
		w.SetContent(newWindow)
		w.Show()

	})
	button.Hidden = true
	hbox := container.NewHBox(button, label)

	list := widget.NewList(
		func() int {
			return len(data.Names)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Test")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data.Names[i])
		})

	// selectedItem := widget.NewLabel()

	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(fmt.Sprintf("Name: %v", data.Names[id]))
		//icon.SetResource(theme.InfoIcon())
		button.Hidden = false
		button.SetIcon(theme.DownloadIcon())
		nameHQ = data.Names[id]
		linkHQ = data.Links[id]
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Selecione uma HQ!")
		//icon.SetResource(nil)
		button.SetIcon(nil)
		button.Hidden = true
		nameHQ = ""
		linkHQ = ""
		labelResult.Text = ""
	}

	input2Search := widget.NewEntry()

	button2Search := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		list.UnselectAll()

		data = utils.HQs{
			Links: []string{},
			Names: []string{},
		}

		data = utils.Search2HQ(input2Search.Text)

		label.Text = "Selecione uma HQ!"

		label.Refresh()

	})

	block1 := container.NewGridWithColumns(2, input2Search, button2Search)

	block2 := container.NewGridWithRows(2, block1, list)

	listView := container.NewHSplit(block2, container.NewCenter(hbox))

	myWindow.SetContent(listView)
	myWindow.ShowAndRun()
}
