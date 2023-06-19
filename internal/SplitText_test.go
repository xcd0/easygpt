package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestSplitText(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	path := "C:/msys64/home/y-hayasaki/work/easygpt/tmp/in/15/15.5._Data_handling.md"
	RunOnBash("cat %v | grep '^#' | wc -l", path)

	inputText := GetTextNoError(&path)
	log.Printf("len(*inputText) : %v", len(*inputText))

	div_lf := DivideOnNewLine(inputText)
	log.Printf("div_lf[last] : %v", div_lf[0][len(div_lf[0])-1])

	headings := DivideOnHeadingElement(div_lf)
	log.Printf("headings[last] : %v", headings[len(headings)-1][len(headings[len(headings)-1])-1])

	divided := DivideBySize(headings, 1)
	log.Printf("divided[last] : %v", divided[len(divided)-1][len(divided[len(divided)-1])-1])

	////////////////////////////////////////////////////////////////////////////////////////////////////
	strs := SplitText(inputText, 1)
	////////////////////////////////////////////////////////////////////////////////////////////////////

	log.Printf("len(strs) : %v", len(strs))

	dname, err := os.MkdirTemp("", "tmp_for_split_text")
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("tmp :", dname)

	for i, str := range strs {
		if len(str) > 15 {
			log.Printf("%v: '%s ... %s'", i, str[:15], str[len(str)-15:])
		} else {
			log.Printf("%v: '%s'", i, str)
		}
		p := filepath.Join(dname, filepath.Base(path)+"_"+strconv.Itoa(i)+".md")
		OutputTextForCheck(&p, &str)
	}
}
