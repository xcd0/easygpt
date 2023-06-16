package internal

import (
	"log"
	"sort"
	"unicode/utf8"

	//"github.com/cloudengio/go.pkgs/algo/codec"
	//"github.com/cloudengio/go.pkgs/algo/lcs"

	"cloudeng.io/algo/codec"
	"cloudeng.io/algo/lcs"
)

func GetCommonParentPath(files []string) string {

	sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })

	for _, f := range files {
		log.Printf("%v", f)
	}

	runeDecoder := codec.NewDecoder(utf8.DecodeRune)
	for i, max := 0, len(files)-1; i < max; i++ {
		a := runeDecoder.Decode([]byte(files[i]))
		b := runeDecoder.Decode([]byte(files[i+1]))
		all := lcs.NewDP(a, b).AllLCS()
		for _, lcs := range all {
			log.Printf("lcs:%v", string(lcs))
		}
	}
	return ""
}
