package internal

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDirwalk(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	path := "../internal"
	apath, _ := filepath.Abs(path)
	t.Logf("Dir: %v", apath)

	ansstr := RunOnBash("tree -d %v | grep 'directories' | awk '{print $1}'", apath)
	t.Logf("ansstr: %v", ansstr)

	ansstr = strings.ReplaceAll(ansstr, "\n", "")
	ans, err := strconv.Atoi(ansstr)
	if err != nil {
		t.Logf("Error: Atoi: %v", err)
	}
	t.Logf("ans : %v", ans)

	t.Logf("tree -d %v | grep 'directories' | awk '{print $1}' -> %v", apath, ans)

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
