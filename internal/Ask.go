package internal

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func Ask(input_file, input_dir, output_dir, tmp_dir, apikey string, prompt *string) {
	// このファイルのファイルパスから固有の文字列を生成し、一時ディレクトリの名前とする。
	id := strings.ReplaceAll(filepath.ToSlash(strings.Replace(input_file, input_dir, "", 1)), "/", "-")[1:]
	// 一時ディレクトリにこのファイル固有のディレクトリを作成する
	id = filepath.Join(tmp_dir, id)
	if len(id) == 0 {
		// なぜか文字がすべて消えた。ファイル名が\とかだったらありえる。
		// ランダム文字列生成してディレクトリ名にする
		uuidObj, _ := uuid.NewUUID()
		id = filepath.Join(tmp_dir, uuidObj.String())
	}
	CreateDir(id)
	OutputTextForCheck(filepath.Join(id, "original_path.txt"), input_file)

	// 出力
	outputpath := strings.Replace(input_file, input_dir, output_dir, 1)

	output := Question(apikey, prompt, input_file, id, true)

	OutputTextForCheck(outputpath, output)
}
