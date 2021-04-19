package main

import (
	"fmt"
	"image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
	"github.com/ribeiroferreiralucas/wpp-scrapper/wppscrapperimp"
)

var application = app.NewWithID("br.ufrrj.wppscrappergui")
var wdw = application.NewWindow("WppScrapper GUI")

var wppScrapper = wppscrapperimp.InitializeConnection().(wppscrapper.IWppScrapper)

func main() {
	application.SetIcon(theme.FyneLogo())

	showInitialLoadingView()

	wdw.Resize(fyne.NewSize(640, 460))
	go initializeWppScrapper()
	wdw.ShowAndRun()
}

func showInitialLoadingView() fyne.CanvasObject {

	title := canvas.NewText("Loaging", theme.TextColor())
	cont := container.New(layout.NewCenterLayout(), title)

	wdw.SetContent(cont)
	return cont
}

func showQrCodeView(qrCode string) fyne.CanvasObject {

	title := canvas.NewText("Scan the QRCode using your WhatsApp app", theme.ForegroundColor())
	qrImage := getQrCodeImage(qrCode)
	qrImage.FillMode = canvas.ImageFillContain
	qrImage.Refresh()
	cont := container.NewBorder(title, nil, nil, nil, qrImage)

	wdw.SetContent(cont)
	return cont
}

func showMainView() {
	mainView := &MainView{}
	// mainView.WppScrapper = wppScrapper
	// mainView.Wdw = wdw

	mainView.Show()
}

func initializeWppScrapper() {

	qrCode := make(chan string)
	go func() {
		showQrCodeView(<-qrCode)
	}()

	_, err := wppScrapper.ReAuth(qrCode, application.UniqueID())
	if err != nil {
		//TODO: tratar erro
		fmt.Println("Error trying to auth", err)

		initializeWppScrapper()
		return
	}

	showInitialLoadingView()

	if !wppScrapper.Initialized() {
		<-wppScrapper.WaitInitialization()
	}

	showMainView()
}

func getQrCodeImage(qrCodeValue string) *canvas.Image {
	qrCode, _ := qr.Encode(qrCodeValue, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	file, _ := os.CreateTemp("", "qrcode_*.png")

	defer file.Close()
	fmt.Println(file.Name())
	// encode the barcode as png
	png.Encode(file, qrCode)

	image := &canvas.Image{
		File: file.Name(),
	}

	return image
}
