package internal

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetSettingFilePath() (*string, error) {
	dirs := []string{GetCurrentDir(), GetHomeDir(), GetBinDir()}
	settingFileNames := []string{
		"easygpt.hjson",
		".easygpt.hjson",
		".easygpt",
	}
	// 設定ファイルは、カレントディレクトリ -> ホームディレクトリ -> 実行ファイルがあるディレクトリ
	// の順で調べる。
	for _, d := range dirs {
		for _, n := range settingFileNames {
			p := filepath.Join(d, n)
			//log.Printf("Debug: check path : %v", p)
			if _, err := os.Stat(p); err == nil {
				// 設定ファイルがあった。
				return &p, nil
			}
		}
	}
	err := errors.Errorf("設定ファイルが見つかりませんでした。")
	//log.Printf("Debug: %v", err)
	return nil, err
}
