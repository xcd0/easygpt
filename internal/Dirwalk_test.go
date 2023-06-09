package internal

import (
	"log"
	"path/filepath"
	"strconv"
	"testing"
)

func TestDirwalk(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	path := "../internal"
	apath, _ := filepath.Abs(path)
	t.Logf("Dir: %v", apath)
	ans := 4
	if true {
		ans, _ = strconv.Atoi(RunOnBash("tree -d %v | grep 'directories' | awk '{print $1}'", apath))
		t.Logf("find %v -type d | wc -l --> %v", apath, ans)
	}

	fes, err := Dirwalk(apath)

	if err != nil {
		for i, d := range fes.Dirs {
			t.Logf("%v: %v", i, d)
		}
		t.Errorf("Error: %v", err)
	}
	dirnum := len(fes.Dirs)
	if dirnum != ans {
		t.Logf("Dirs: %v", fes.Dirs)
		t.Errorf("get: %v, expect: %v", dirnum, ans)
	}
}
