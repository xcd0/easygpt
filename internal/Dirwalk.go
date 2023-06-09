package internal

import (
	"os"
	"path/filepath"
)

type FileEntries struct {
	Dirs  []string
	Files []string
}

// find -type dとは違って、指定したディレクトリを含まない
// treeと同じカウント方法
// $ tree -d --charset ascii ./input | grep "directories"
// 4 directories
func Dirwalk(dirpath string) (FileEntries, error) {
	var dirs []string
	var files []string

	if f, err := os.Stat(dirpath); err != nil {
		return FileEntries{}, err // 不明
	} else if f.IsDir() {
		// ここですると一番親のディレクトリも入る
		//dirs = append(dirs, dirpath)
	} else {
		return FileEntries{}, nil // 不明
	}

	dirpath, _ = filepath.Abs(dirpath)
	dirpath = filepath.Clean(dirpath)

	direntry, err := os.ReadDir(dirpath)
	if err != nil {
		return FileEntries{}, err
	}
	for _, e := range direntry {
		p := filepath.Join(dirpath, e.Name())
		p = filepath.Clean(p)
		if e.IsDir() {
			//log.Printf("Dir : %v", p)
			// Dirwalkの先頭でappendされる。
			// やめた。
			dirs = append(dirs, p)
			fe, err := Dirwalk(p)
			//
			if err != nil {
				return FileEntries{}, err
			}
			files = append(files, fe.Files...)
			dirs = append(dirs, fe.Dirs...)
		} else {
			//log.Printf("File: %v", p)
			files = append(files, p)
		}
	}
	return FileEntries{Dirs: dirs, Files: files}, nil
}
