package internal

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
)

/*
D&Dによる実行を想定した処理。

* どこに出力させるか
* 設定ファイルのoutput-dirに設定が
	* ある場合
		* 入力ファイル群のパスが
			* /hoge/a/b.txt
			* /hoge/a/b/c.txt
			* /hoge/a/b/c/d.txt
		のように/a如何にあるとわかれば、output-dirにaを作成して保存したい。
	* ない場合
		* 入力ファイルと同じディレクトリに名前を変えて出力。

*/

func RunByDD(settingsfordd []SettingForDD, files []string, setting *Setting) {

	common := ""
	if len(settingsfordd[0].OutputDir) != 0 {
		// 与えられたファイルパス群に共通の親ディレクトリを探す。
		common = GetCommonParentPath(files)
		// commonのサブディレクトリと同じディレクトリ構造を、出力先ディレクトリに作成する。
		CreateSameDirsOn(common, settingsfordd[0].OutputDir)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, setting.Concurrency) // concurrency数のバッファ
	pbs := make([]*pb.ProgressBar, len(files))
	for i, f := range files {
		pbs[i] = pb.New(2).Prefix(fmt.Sprintf("%-50s", filepath.Base(f)))
	}
	pool, err := pb.StartPool(pbs...)
	if err != nil {
		//panic(err) 無視
	}

	for i := 0; i < len(settingsfordd); i++ {
		s := settingsfordd[i]
		p := pbs[i]
		sem <- struct{}{}
		wg.Add(1)
		go func(s *SettingForDD, common string, p *pb.ProgressBar) {
			defer wg.Done()
			defer func() { <-sem }() // 処理が終わったらチャネルを解放
			defer p.Finish()
			AskByDD(s, common, setting)
		}(&s, common, p)
	}
	wg.Wait()
	pool.Stop()
}

func AskByDD(s *SettingForDD, lcs string, setting *Setting) {

	outputpath := func() string {
		if len(s.OutputDir) != 0 {
			// 出力先ディレクトリの設定がある。
			ofname, _ := filepath.Abs(s.Input)               // 絶対パスに変換
			ofname = strings.Replace(ofname, lcs, "", 1)     // 共通部分を削除
			ofname = *AddPostfix(&ofname, &s.Common.Postfix) // postfixを付与
			return filepath.Join(s.OutputDir, ofname)        // lcsが/a/b/のとき、上記の処理で/a/b/c/dのofnameはc/d_postfixになり、それを出力ディレクトリのパスに付与する。
		} else {
			// 出力先ディレクトリの設定がない
			// 入力ファイルと同じディレクトリに、名前を変えて出力する。
			if len(s.Common.Postfix) == 0 {
				// postfixが空の時、出力ファイル名が入力ファイルと同じになってしまい、上書きしてしまう。
				fmt.Printf("設定ファイルにおいて、postfixが設定されていません。")
				fmt.Printf("")
			}
			ofname := AddPostfix(&s.Input, &s.Common.Postfix) // postfixを付与
			return filepath.Join(filepath.Dir(s.Input), *ofname)
		}
	}()

	inputText := GetTextNoError(&s.Input)
	outputText, err := QuestionByText(inputText, setting, false, nil)
	if err != nil {
		log.Printf("Error: %v, %v", s.Input, err)
	} else {
		OutputTextForCheck(&outputpath, outputText)
	}
	fmt.Printf("end   : %v\n", s.Input)
}

func AddPostfix(file, postfix *string) *string {

	// これはすでに最低限必要な設定が得られているときにのみ生成する。最低限の設定が得られていないときnilになる。
	// postfixが空の時、入力ファイルに上書きしてしまうので勝手に文字列を付与する。
	if postfix == nil {
		return file
	}
	if len(*postfix) == 0 {
		*postfix = "_easygpt_output"
	}
	// どうしても上書きしたいときは `--postfix *` のように設定する。
	if *postfix == "*" {
		*postfix = ""
	}
	ext := filepath.Ext(*file)
	name := GetFileNameWithoutExt(file)
	n := fmt.Sprintf("%v%v%v", name, *postfix, ext)
	return &n
}
