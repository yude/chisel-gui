package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	// 環境を設定
	os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\meiryo.ttc")
	exec.Command("chcp 65001")
	// GUI を設定
	a := app.New()
	w := a.NewWindow("chisel - A fast TCP tunnel over HTTP")

	hello, desc := widget.NewLabelWithStyle("Minecraft 接続用プロキシ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), widget.NewLabelWithStyle("chisel.exe が存在しているか確かめます。", fyne.TextAlignCenter, fyne.TextStyle{})
	var launch *widget.Button
	logField := widget.NewMultiLineEntry()
	// 実行ログ関係
	logField.SetPlaceHolder("ここに実行ログが表示されます。")
	logField.Disable()
	// ボタンの動画
	launch = widget.NewButton("起動する", func() {
		// chisel のダウンロード元をOSによって分岐
		var url string
		if runtime.GOOS == "windows" {
			url = "https://github.com/jpillora/chisel/releases/download/v1.6.0/chisel_1.6.0_windows_amd64.gz"
		}
		if runtime.GOOS == "darwin" {
			url = "https://github.com/jpillora/chisel/releases/download/v1.6.0/chisel_1.6.0_darwin_amd64.gz"
		}
		var newfilename string
		if runtime.GOOS == "windows" {
			newfilename = "chisel.exe"
		}
		if runtime.GOOS == "darwin" {
			newfilename = "chisel"
		}
		// ボタンの状態を変更
		launch.SetText("起動済み")
		launch.Disable()
		// chisel をダウンロードして展開する
		if f, err := os.Stat(newfilename); os.IsNotExist(err) || f.IsDir() {
			// chisel (gz ファイル) をダウンロードする
			desc.SetText("chisel が存在しません。ダウンロードします。")
			output, err := os.Create("chisel.gz")
			defer output.Close()
			response, err := http.Get(url)
			desc.SetText("chisel が存在しません。ダウンロードしています...")
			if err != nil {
				desc.SetText("chisel を何らかの理由でダウンロードできませんでした。")
				return
			}
			defer response.Body.Close()
			n, err := io.Copy(output, response.Body)
			desc.SetText(fmt.Sprintf("%d バイトをダウンロードしました。展開します。", n))
			time.Sleep(1 * time.Second)
			// 展開部
			desc.SetText("展開しています...")
			filename := "chisel.gz"
			gzipfile, err := os.Open(filename)

			if err != nil {
				logField.SetText(fmt.Sprintln(err))
				println(err)
				//os.Exit(1)
			}

			reader, err := gzip.NewReader(gzipfile)
			if err != nil {
				logField.SetText(fmt.Sprintln(err))
				println(err)
				//os.Exit(1)
			}
			defer reader.Close()

			writer, err := os.Create(newfilename)

			if err != nil {
				logField.SetText(fmt.Sprintln(err))
				println(err)
				// os.Exit(1)
			}

			defer writer.Close()
			if _, err = io.Copy(writer, reader); err != nil {
				logField.SetText(fmt.Sprintln(err))
				println(err)
				// os.Exit(1)
			}
			defer gzipfile.Close()
			// ダウンロードした chisel.gz を削除
			if err := os.Remove("chisel.gz"); err != nil {
				fmt.Println(err)
			}
			desc.SetText("展開しました。プロキシを起動します。")
		} else {
			desc.SetText("chisel が既に存在します。プロキシを起動します。")
		}
		time.Sleep(1 * time.Second)
		// chisel をバックグラウンドで起動
		c := "chisel.exe"
		p := []string{"client", "proxy.yude.moe", "25565"}
		cmd := exec.Command(c, p...)
		cmd.Dir = "."
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

		stdoutScanner := bufio.NewScanner(stdout)
		stdoutOutputChan := make(chan string)
		stdoutDoneChan := make(chan bool)
		stderrScanner := bufio.NewScanner(stderr)
		stderrOutputChan := make(chan string)
		stderrDoneChan := make(chan bool)
		go streamReader(stdoutScanner, stdoutOutputChan, stdoutDoneChan)
		go streamReader(stderrScanner, stderrOutputChan, stderrDoneChan)

		stillGoing := true
		for stillGoing {
			select {
			case <-stdoutDoneChan:
				stillGoing = false
			case line := <-stdoutOutputChan:
				// log.Println(line)
				logField.SetText(line)
			case line := <-stderrOutputChan:
				// log.Println(line)
				logField.SetText(line)
				/*if fmt.Println(strings.LastIndex(line, "test")) {
					println("334")
				}*/
			}
		}

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
