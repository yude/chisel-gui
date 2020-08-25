package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {

	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("#いちぴろマイクラ部 接続用プロキシ\n")
	fmt.Printf("Powered by chisel, A fast TCP tunnel over HTTP (https://github.com/jpillora/chisel)\n")
	fmt.Printf("\n")
	if f, err := os.Stat("chisel.exe"); os.IsNotExist(err) || f.IsDir() {
		url := "https://chiselmirror.now.sh/chisel.exe"
		fmt.Println("chisel が存在しません。ダウンロードします。")
		output, err := os.Create("chisel.exe")
		defer output.Close()

		response, err := http.Get(url)
		if err != nil {
			fmt.Println("chisel を何らかの理由でダウンロードできませんでした。", url, "-", err)
			return
		}
		defer response.Body.Close()

		n, err := io.Copy(output, response.Body)

		fmt.Println(n, "バイトをダウンロードしました。")
	} else {
		fmt.Println("chisel が見つかりました。プロキシを起動します。")
	}
		cmdName := "chisel.exe client proxy.yude.moe 25565"
		cmdArgs := strings.Fields(cmdName)

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		oneByte := make([]byte, 100)
		num := 1
		for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		fmt.Println(string(line))
		num = num + 1
		if num > 3 {
			fmt.Printf("長期間サーバーから切断されたため、自動的にプロキシを閉じました。\n")
			fmt.Printf("再度接続する場合、もう一度 chisel_proxy.exe を実行してください。\n")
			fmt.Printf("終了するにはこのウィンドウを閉じてください.")
			time.Sleep(time.Second * 10)
		}
	}