package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\meiryo.ttc")
	a := app.New()
	w := a.NewWindow("chisel - A fast TCP tunnel over HTTP")
	hello, desc := widget.NewLabelWithStyle("こんにちは！", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), widget.NewLabelWithStyle("chisel.exe が存在しているか確かめます。", fyne.TextAlignCenter, fyne.TextStyle{})
	launch := widget.NewButton("起動する", func() {
		if f, err := os.Stat("chisel.exe"); os.IsNotExist(err) || f.IsDir() {
			desc.SetText("chisel.exe が存在しません。ダウンロードします。")
			url := "https://chiselmirror.now.sh/chisel.exe"
			output, err := os.Create("chisel.exe")
			defer output.Close()

			response, err := http.Get(url)
			desc.SetText("chisel.exe が存在しません。ダウンロードしています...")
			if err != nil {
				desc.SetText("chisel.exe を何らかの理由でダウンロードできませんでした。")
				return
			}
			defer response.Body.Close()
			n, err := io.Copy(output, response.Body)
			desc.SetText(fmt.Sprintf("%d バイトをダウンロードしました。", n))
		} else {
			desc.SetText("chisel.exe が存在します。プロキシを起動します。")
		}
	})
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(390, 100))
	w.SetContent(widget.NewVBox(
		hello,
		desc,
		launch,
	))
	w.ShowAndRun()

}
