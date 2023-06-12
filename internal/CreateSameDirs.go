package internal

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func AbsPath(path *string) {
	*path, _ = filepath.Abs(*path)
	*path = filepath.Clean(*path)
}

func CreateSameDirsOn(target_dir string, ref_dir string) {

	//fmt.Printf("%+v", errors.Errorf("debug"))

	// ref_dirの子ディレクトリをtarget_dirに作成する。
	// ref_dirの子ディレクトリ一覧からref_dirのパス文字列を削除して、target_dirをくっつける
	AbsPath(&target_dir)
	AbsPath(&ref_dir)

	//log.Printf("target: %v, ref_dir: %v", target_dir, ref_dir)
	de, err := Dirwalk(ref_dir)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//log.Printf("dirs: %v", de)

	// 出力先ディレクトリを作っておく
	if err := os.Mkdir(target_dir, 0777); err != nil {
		fe, e := "file exists", err.Error()
		if ee := e[len(e)-len(fe):]; ee != fe {
			log.Printf("Error: %v", err)
		}
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, d := range de.Dirs {
		rep := strings.Replace(d, ref_dir, "", 1) // 元のパスを削除する
		d = filepath.Join(target_dir, rep)
		//log.Printf("%v ", d)
		if err := os.Mkdir(d, 0777); err != nil {
			fe, e := "file exists", err.Error()
			if ee := e[len(e)-len(fe):]; ee != fe {
				log.Printf("Error: %v", err)
			}
		}
	}
	return
}
