package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
)

type MainView struct {
	cnvProgressBar     *widget.ProgressBar
	lblStatus          *widget.Label
	btnStartScrap      *widget.Button
	btnRestartScrap    *widget.Button
	btnStopScrap       *widget.Button
	header             *fyne.Container
	scrappedChatsCount int
}

func (m *MainView) Show() {

	eventHandler := wppScrapper.GetWppScrapperEventHandler()
	eventHandler.AddOnChatScrapFinishedListener(m)
	eventHandler.AddOnChatScrapStartedListener(m)
	eventHandler.AddOnScrapperFinishedListener(m)
	eventHandler.AddOnScrapperStartedListener(m)
	eventHandler.AddOnScrapperStoppedListener(m)

	m.buildView()
}

func (m *MainView) OnWppScrapperStarted(wppScrapper wppscrapper.IWppScrapper) {
	m.updateButtons(true)
	m.lblStatus.SetText("Status: Running")
}

func (m *MainView) OnWppScrapperStopped(wppScrapper wppscrapper.IWppScrapper) {
	m.updateButtons(false)
	m.lblStatus.SetText("Status: Stopped")
}
func (m *MainView) OnWppScrapperFinished(wppScrapper wppscrapper.IWppScrapper) {
	m.updateButtons(false)
	m.lblStatus.SetText("Status: Finished")
}
func (m *MainView) OnWppScrapperChatScrapStarted(chat wppscrapper.Chat) {
	wdw.Content().Refresh()
}
func (m *MainView) OnWppScrapperChatScrapFinished(chat wppscrapper.Chat) {

	m.scrappedChatsCount++
	progress := float64(m.scrappedChatsCount) / float64(len(wppScrapper.GetChats()))
	m.cnvProgressBar.SetValue(progress)
	wdw.Content().Refresh()
}

func (m *MainView) buildView() {
	header := m.buildHeader()
	table := m.buildChatsTable()
	progressBar := m.buildProgressBar()

	cont := container.NewBorder(header, progressBar, nil, nil, table)

	wdw.SetContent(cont)
}

func (m *MainView) buildProgressBar() fyne.CanvasObject {

	m.cnvProgressBar = widget.NewProgressBar()
	return m.cnvProgressBar
}

func (m *MainView) buildHeader() *fyne.Container {

	m.lblStatus = &widget.Label{
		Text: "Status: Idle",
	}

	m.btnRestartScrap = &widget.Button{
		Text: "Start Scrapper",
		Icon: theme.DownloadIcon(),
		OnTapped: func() {
			m.scrappedChatsCount = 0
			wppScrapper.StartScrapper(false)
		},
	}
	m.btnStartScrap = &widget.Button{
		Text: "Restart Scrapper",
		Icon: theme.DownloadIcon(),
		OnTapped: func() {
			wppScrapper.StartScrapper(true)
		},
	}

	m.btnStopScrap = &widget.Button{
		Text: "Stop Scrapper",
		Icon: theme.CancelIcon(),
		OnTapped: func() {
			wppScrapper.StopScrapper()
		},
	}
	m.btnStopScrap.Hide()
	m.header = container.NewHBox(m.lblStatus, layout.NewSpacer(), m.btnStartScrap, m.btnRestartScrap, m.btnStopScrap)
	m.header.Refresh()
	return m.header
}

func (m *MainView) buildChatsTable() fyne.CanvasObject {
	chats := wppScrapper.GetChats()

	tableValues := []string{}
	for key, _ := range chats {
		tableValues = append(tableValues, key)
	}

	size := len(chats)

	resultantWidht := [3]int{0, 0, 0}

	table := &widget.Table{
		Length: func() (int, int) {

			return size, 3
		},
		CreateCell: func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Resize(fyne.NewSize(100, 50))
			return label

		},
	}
	table.UpdateCell = func(id widget.TableCellID, template fyne.CanvasObject) {
		label := template.(*widget.Label)
		if id.Col == 0 {
			chats := wppScrapper.GetChats()
			chatId := tableValues[id.Row]
			statusDesc := statusToStrng(chats[chatId].GetStatus())
			label.SetText(statusDesc)
		}
		if id.Col == 1 {
			idLabel := tableValues[id.Row]

			label.SetText(idLabel)

		}
		if id.Col == 2 {
			chats := wppScrapper.GetChats()
			chatId := tableValues[id.Row]
			label.SetText(chats[chatId].Name())
			label.Refresh()
		}

		if (10 * len(label.Text)) > resultantWidht[id.Col] {
			resultantWidht[id.Col] = 10 * len(label.Text)
		}
		label.Refresh()

	}
	table.SetColumnWidth(0, 150)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 300)
	return table
}

func (m *MainView) updateButtons(isRunning bool) {

	if isRunning {
		m.btnStopScrap.Show()
		m.btnStartScrap.Hide()
		m.btnRestartScrap.Hide()
	} else {
		m.btnStopScrap.Hide()
		m.btnStartScrap.Show()
		m.btnRestartScrap.Show()
	}
	m.header.Refresh()
}

func statusToStrng(status wppscrapper.ChatStatus) string {
	switch status {
	case wppscrapper.Idle:
		return "Idle"
	case wppscrapper.Queue:
		return "Queue"
	case wppscrapper.Running:
		return "Running"
	case wppscrapper.Stoped:
		return "Stoped"
	case wppscrapper.Finished:
		return "Finished"
	}
	return "Unknown Status"
}
