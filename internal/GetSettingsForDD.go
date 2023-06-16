package internal

import (
	"path/filepath"
)

func GetSettingsForDD(argsAll *ArgsAll, setting *Setting) ([]SettingForDD, []string) {
	// D&Dで実行された場合の対応のために、 引数に--xxxのような形式でない引数があれば、 それは個別に実行する。
	if len(argsAll.InputFiles) == 0 {
		return nil, nil // ポジショナル引数がなかった。別によい。
	}
	// ポジショナル引数があった。引数一つ毎に設定を生成する。
	// ファイルはファイルごとに処理する。ディレクトリは再帰的にファイルを探索して処理する。
	files := GetFileList(argsAll.InputFiles)
	var settingsfordd []SettingForDD

	for _, f := range files {
		settingsfordd = append(settingsfordd,
			SettingForDD{
				Input:     f,                 // 入力ファイル
				OutputDir: setting.OutputDir, // 入力ファイル
				Common:    &setting.SettingCommon,
			},
		)
	}
	return settingsfordd, files

}

func GetFileList(paths []string) []string {
	// ファイル一覧を生成
	var files []string
	for _, p := range paths {
		p, _ = filepath.Abs(p)
		CheckPathAndRunFunction(p,
			func(err error) {}, // 存在しない時無視する。
			func() { // ディレクトリの場合
				fes, err := Dirwalk(p)
				if err != nil { // エラーの場合でも無視する。
					//fmt.Printf("Error: %v", err)
				}
				files = append(files, fes.Files...)
			},
			func() { // ファイルが存在する場合
				files = append(files, p)
			},
		)
	}
	return files
}
