package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\meiryo.ttc")
	exec.Command("chcp 65001")
	a := app.New()
	w := a.NewWindow("chisel - A fast TCP tunnel over HTTP")
	hello, desc := widget.NewLabelWithStyle("Minecraft 接続用プロキシ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), widget.NewLabelWithStyle("chisel.exe が存在しているか確かめます。", fyne.TextAlignCenter, fyne.TextStyle{})
	var launch *widget.Button
	logField := widget.NewMultiLineEntry()
	logField.SetPlaceHolder("ここに実行ログが表示されます。")
	logField.Disable()
	launch = widget.NewButton("起動する", func() {
		launch.SetText("起動済み")
		launch.Disable()

		if f, err := os.Stat("chisel_x64.exe"); os.IsNotExist(err) || f.IsDir() {
			desc.SetText("chisel.exe が存在しません。ダウンロードします。")
			url := "https://chiselmirror.now.sh/chisel_x64.exe"
			output, err := os.Create("chisel")
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
			time.Sleep(1 * time.Second)
			desc.SetText("プロキシを起動します。")
		} else {
			desc.SetText("chisel.exe が既に存在します。プロキシを起動します。")
		}
		// コマンドと引数を定義する
		c := "chisel_x64.exe"
		p := []string{"client", "proxy.yude.moe", "25565"}
		cmd := exec.Command(c, p...)

		// 実行ディレクトリ
		cmd.Dir = "."

		// パイプを作る
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		streamReader := func(scanner *bufio.Scanner, outputChan chan string, doneChan chan bool) {
			defer close(outputChan)
			defer close(doneChan)
			for scanner.Scan() {
				outputChan <- scanner.Text()
			}
			doneChan <- true
		}

		// stdout, stderrをひろうgoroutineを起動
		stdoutScanner := bufio.NewScanner(stdout)
		stdoutOutputChan := make(chan string)
		stdoutDoneChan := make(chan bool)
		stderrScanner := bufio.NewScanner(stderr)
		stderrOutputChan := make(chan string)
		stderrDoneChan := make(chan bool)
		go streamReader(stdoutScanner, stdoutOutputChan, stdoutDoneChan)
		go streamReader(stderrScanner, stderrOutputChan, stderrDoneChan)

		// channel経由でデータを引っこ抜く
		stillGoing := true
		for stillGoing {
			select {
			case <-stdoutDoneChan:
				stillGoing = false
			case line := <-stdoutOutputChan:
				log.Println(line)
				logField.SetText(line)
			case line := <-stderrOutputChan:
				log.Println(line)
				logField.SetText(line)
			}
		}

		//一応Waitでプロセスの終了をまつ
		ret := cmd.Wait()
		if ret != nil {
			log.Fatal(err)
		}
	})
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(390, 100))
	w.SetContent(widget.NewVBox(
		hello,
		desc,
		launch,
		logField,
	))
	w.ShowAndRun()

}
