package internal

import (
	"fmt"
	"path/filepath"
)

func GenerateSettingForDD(apikey, prompt, postfix *string, inputArgFiles *[]string) ([]SettingForDD, []string) {
	// D&Dによる実行
	// D&Dの場合、引数を舐めて、ファイルはファイルごとに実行
	// ディレクトリは設定ありとある程度同じように実行

	// ファイル一覧を生成
	var files []string
	for _, p := range *inputArgFiles {
		p, _ = filepath.Abs(p)
		CheckPathAndRunFunction(p,
			func(err error) {
				// 何もしなければcontinue相当の動作になる。
				// 存在しない時無視する。
			},
			func() {
				fes, err := Dirwalk(p)
				if err != nil {
					fmt.Printf("Error: %v", err)
					fmt.Printf("       Skipped.")
				}
				files = append(files, fes.Files...)
			},
			func() {
				// ファイルが存在する
				files = append(files, p)
			},
		)
	}
	var settings []SettingForDD
	for _, f := range files {
		settings = append(settings,
			SettingForDD{
				ApiKey:      *apikey,  // APIキー
				Input:       f,        // 入力ファイル
				Prompt:      *prompt,  // プロンプト
				Postfix:     *postfix, // 出力ファイル名の末尾に付与する文字列。空の時 "_easygpt_output" となる。
				Tmp:         "",       // 一時ファイルを保存するディレクトリ
				Concurrency: 1,        // 最大同時並列実行数
			},
		)
	}

	return settings, files
}
