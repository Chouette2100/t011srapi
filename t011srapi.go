/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

package main

import (
	"io"
	"log"
	"os"
	"time"

	"net/url"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

/*

	星・種を集め、星集め・種集めの状況を調べる

	ソースのダウンロード、ビルドについて以下簡単に説明します。詳細は以下の記事を参考にしてください。
	WindowsもLinuxも特記した部分以外は同じです。

		【Windows】かんたんなWebサーバーの作り方
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/c5cab5

		---------------------

		【Windows】Githubにあるサンプルプログラムの実行方法
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/e27fc9

		【Unix/Linux】Githubにあるサンプルプログラムの実行方法
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/220e38

			ロードモジュールさえできればいいということでしたらコマンド一つでできます。

【Unix/Linux】

	$ cd ~/go/src
	$ curl -OL https://github.com/Chouette2100/t011srapi/archive/refs/tags/v0.1.0.tar.gz
	$ tar xvf v0.1.0.tar.gz
	$ mv t011srapi-0.1.0 t011srapi
	$ cd t011srapi
	$ go mod init
	$ go mod tidy
	$ go build t011srapi.go
	
	実行する（下記は例）

	$ ./t011srapi 48_Kaori_Inagaki


	【Windows】

	Microsoft Windows [Version 10.0.22000.856]
	(c) Microsoft Corporation. All rights reserved.

	C:\Users\chouette>cd go

	C:\Users\chouette\go>cd src

	作業はかならず %HOMEPATH%\go\src の下で行います。

	以下、要するに https://github.com/Chouette2100/t011srapi/releases にあるv0.1.0のZIPファイルSource code (zip) からソースをとりだしてくださいということなので、ブラウザでダウンロードしてエクスプローラで解凍というこでもけっこうです。なんならこの記事の最後にあるgithubのソースをエディターにコピペで作るということでもかまいません（この場合文字コードはかならずUTF-8にしてください 改行はLFになっています。というようなことを考えるとやっぱりダウンロードして解凍が安全かも）

	C:\Users\chouette\go\src>mkdir t011srapi

	C:\Users\chouette\go\src>cd t011srapi

	C:\Users\chouette\go\src\t011srapi>curl -OL https://github.com/Chouette2100/t011srapi/archive/refs/tags/v0.1.0.zip
	  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
	                                 Dload  Upload   Total   Spent    Left  Speed
	  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
	100  6265    0  6265    0     0   6777      0 --:--:-- --:--:-- --:--:-- 16400

	C:\Users\chouette\go\src\t011srapi>call powershell -command "Expand-Archive v0.1.0.zip"

	C:\Users\chouette\go\src\t011srapi>tree
	フォルダー パスの一覧
	ボリューム シリアル番号は E2CD-BDF1 です
	C:.
	└─v0.1.0
	    └─t011srapi-0.1.0
	        ├─public
	        └─templates

	C:\Users\chouette\go\src\t011srapi>xcopy /e v0.1.0\t011srapi-0.1.0\*.*
	v0.1.0\t011srapi-0.1.0\freebsd.bat
	v0.1.0\t011srapi-0.1.0\freebsd.sh
	v0.1.0\t011srapi-0.1.0\LICENSE
	v0.1.0\t011srapi-0.1.0\README.md
	v0.1.0\t011srapi-0.1.0\t011srapi.go
	v0.1.0\t011srapi-0.1.0\public\index.html
	v0.1.0\t011srapi-0.1.0\templates\top.gtpl
	以下省略（リリースによって内容が異なる場合があります。）

	C:\Users\chouette\go\src\t011srapi>rmdir /s /q v0.1.0
	以下省略（リリースによって内容が異なる場合があります。）

	C:\Users\chouette\go\src\t011srapi>del v0.1.0.zip

	ここで次のような構成になっていればOKです。top.gtpl と index.html が所定の場所にあることをかならず確かめてください。

	C:%HOMEPATH%\go\src\t011srapi --+-- t011srapi.go
	                                |
	                                +-- \templates --- t011top.gtpl
	                                |
	                                +-- \public    --- index.html

	ここからはコマンド三つでビルドが完了します。

	C:\Users\chouette\go\src\t011srapi>go mod init
	go: creating new go.mod: module t011srapi
	go: to add module requirements and sums:
	        go mod tidy

	C:\Users\chouette\go\src\t011srapi>go mod tidy
	go: finding module for package github.com/dustin/go-humanize
	go: downloading github.com/dustin/go-humanize v1.0.0
	go: found github.com/dustin/go-humanize in github.com/dustin/go-humanize v1.0.0

	C:\Users\chouette\go\src\t011srapi>go build t011srapi.go

	あとはたとえば次のように実行します

	C:\Users\chouette\go\src\t011srapi>t011srapi 48_Kaori_Inagaki


	Ver. 0.1.0

*/

type Config struct {
	SR_acct string //    SHOWROOMのアカウント名
	SR_pswd string //    SHOWROOMのパスワード
}

func main() {

	//	ログファイルを作る。
	logfilename := time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()

	//	ログをコンソールにも出力する。
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	if len(os.Args) != 2 {
		log.Printf("usage: %s room_url_key\n", os.Args[0])
		os.Exit(1)
	}
	room_url_key := os.Args[1]

	//	設定ファイルの内容を読み込む
	config := new(Config)
	err = exsrapi.LoadConfig("config.yml", config)
	if err != nil {
		log.Printf("exsrapi.LoadConfig(): %s\n", err.Error())
		return
	}

	//	HTTPクライアントとcookiejarを作る
	client, jar, err := exsrapi.CreateNewClient(config.SR_acct)
	if err != nil {
		log.Printf("exsrapi.CreateNewClient(): %s\n", err.Error())
		return
	}
	defer jar.Save()

	//	SHOWROOMにログインした状態にする。
	userid, err := exsrapi.LoginShowroom(client, config.SR_acct, config.SR_pswd)
	if err != nil {
		log.Printf("exsrapi.LoginShowroom(): %s\n", err.Error())
		return
	}
	if userid == "0" {
		log.Printf("ログインできませんでした。アカウントとパスワードを確認してください。\n")
		return
	}

	//	ルーム情報を取得する。
	rs, err := srapi.ApiRoomStatus(client, room_url_key)
	if err != nil {
		log.Printf("srapi.ApiRoomStatus(): %s\n", err.Error())
		return
	}
	if ! rs.Is_live {
		log.Printf("room_url_key = %s　のルームは現在配信が行われていません。\n", room_url_key)
		return
	}
	log.Printf("room_url_key = %s, room_id = %d\n", room_url_key, rs.Room_id)

	//	配信ルームに接続する。
	turl := "https://www.showroom-live.com/" + room_url_key

	u, err := url.Parse(turl)
	if err != nil {
		log.Printf("url.Parse(): %s\n", err.Error())
		return
	}
	resp, err := client.Get(u.String())
	if err != nil {
		log.Printf("http.Get(): %s\n", err.Error())
		return
	}

	//	接続したまま31秒待つ
	time.Sleep(31 * time.Second)
	resp.Body.Close()

	//	星集め、種集めの状況を確認する。
	lp, err := srapi.ApiLivePolling(client, rs.Room_id)
	if err != nil {
		log.Printf("srapi.ApiLivePolling(): %s\n", err.Error())
		return
	}

	log.Printf("OK? %d, IsFree? %s, Toast=[%s]\n", lp.Live_watch_incentive.Ok, lp.Live_watch_incentive.Is_amateur, lp.Toast.Message)

}
