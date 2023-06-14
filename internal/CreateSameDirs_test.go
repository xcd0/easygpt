package internal

import (
	"log"
	"os"
	"testing"
)

func TestCreateSameDirsOn(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	setting := Setting{
		InputDir:  "./for_test/input",
		OutputDir: "./for_test/output",
		Prompt:    "aaa",
	}

	os.RemoveAll(setting.OutputDir)
	if err := os.Mkdir(setting.OutputDir, 0777); err != nil {
		fe, e := "file exists", err.Error()
		if ee := e[len(e)-len(fe):]; ee != fe {
			log.Printf("Error: %v", err)
		}
	}

	CreateSameDirsOn(setting.OutputDir, setting.InputDir)

	ifes, _ := Dirwalk(setting.InputDir)
	ofes, _ := Dirwalk(setting.OutputDir)
	idirnum := len(ifes.Dirs)
	odirnum := len(ofes.Dirs)
	if idirnum == odirnum {
		//
	} else {
		RunOnBash("tree -d --charset ascii %v", setting.InputDir)
		count := RunOnBash("tree -d --charset ascii %v | grep 'directories' | awk '{print $1}'", setting.InputDir)
		t.Errorf("output: %v, expected: %v, tree: %v", odirnum, idirnum, count)
	}
}
