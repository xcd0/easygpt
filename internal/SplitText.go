package internal

import (
	"log"
	"regexp"
	"strings"
)

const (
	Magnification = 1000
)

var (
	reg_elem_h   = regexp.MustCompile("^#")
	reg_elem_pre = regexp.MustCompile("```")
	//reg_elem_empty = regexp.MustCompile("^\n") // 本来これ
	reg_elem_empty         = regexp.MustCompile("\n$") // 一度全て改行で分割した後、空行を前の要素の末尾に\nとして付与しているのでこれで空行判定する。
	reg_elem_func_end      = regexp.MustCompile("^}$")
	reg_elem_comment_start = regexp.MustCompile("^$// ")
)

// var IsOver func(value, num int) bool
func IsOver(value, num int) bool {
	return value >= num*Magnification
}

// sizeに指定された値
func SplitText(text *string, num int) []string {

	*text = strings.ReplaceAll(*text, "\r\n", "\n") // delete CR

	if num <= 0 {
		return []string{*text}
	} else if num > 16 {
		num = 16
	}

	//IsOver = func(value, num int) bool {
	//	return value >= num*Magnification
	//}

	log.Printf("debug: inputText len : %v", len(*text))

	if !IsOver(len(*text), num) {
		return []string{*text} // 区切りサイズより小さかったので分割せず、そのまま返す
	}

	lines := DivideOnNewLine(text)

	//log.Printf("lines: %v", lines[0][0])
	headings := DivideOnHeadingElement(lines)

	// {{{
	//log.Printf("head len: %v", debug_count_slice(headings))
	/*
		for i, strs := range headings {
			for j, str := range strs {
				if len(str) > 10 {
					log.Printf("head: %v:%v: '%s...'", i, j, str[:9])
				} else {
					log.Printf("head: %v:%v: '%s'", i, j, str)
				}
			}
		}

		log.Printf("len(head): %v", len(headings))
		for i := 0; i < 5; i++ {
			if len(headings[i]) > 1 {
				log.Printf("head: %v", headings[i][1])
			} else {
				log.Printf("head: %v", headings[i])
			}
		}
	*/
	// }}}
	divided := DivideBySize(headings, num)
	// {{{
	//log.Printf("div len: %v", debug_count_slice(divided))

	/*
		for i, strs := range divided {
			for j, str := range strs {
				if len(str) > 10 {
					log.Printf("div: %v:%v: '%s...'", i, j, str[:9])
				} else {
					log.Printf("div: %v:%v: '%s'", i, j, str)
				}
			}
		}
	*/
	// }}}
	divided = DivideByFuncEnd(divided, num)

	out := []string{}
	for _, d := range divided {
		out = append(out, strings.Join(d, "\n"))
	}

	return out
}

// dのサイズを計算
func calcStrArrayLength(d []string) int {
	s := 0
	for _, line := range d {
		s += len(line)
	}
	return s
}

func DivideByFuncEnd(divided [][]string, num int) [][]string { // {{{
	out := [][]string{}

	for _, d := range divided {
		// dのサイズを計算
		s := calcStrArrayLength(d)
		//log.Printf("s:%v, num:%v", s, num)
		//log.Printf("-------------------------------------------------------------")
		if !IsOver(s, num) {
			//log.Printf("=============================================================")
			// 越えていないのでそのまま維持
			log.Printf("divideByFuncEnd 越えていないのでそのまま維持 s : %v, num : %v", s, num)
			out = append(out, d)
			continue
		}
		//log.Printf("-------------------------------------------------------------")
		// サイズが超えているので分割
		flagPre := false // <pre>要素内で区切りたくない。
		// 一度bufに入れて、越えたときのbufをoutに吐く
		var buf []string // 越えていない状態の文字列群
		bufsize := 0
		for _, str := range d {
			// 行ごとの処理

			if reg_elem_pre.MatchString(str) { // <pre>に相当する```にマッチした
				flagPre = !flagPre
			}
			if flagPre { // <pre>要素の中の時
				bufsize += len(str)
				buf = append(buf, str)
			} else {
				// <pre>要素でない時
				// 空行があれば、その直前まででサイズを計算

				if !reg_elem_func_end.MatchString(str) { // ^}$にマッチ
					bufsize += len(str)
					buf = append(buf, str)
				} else {
					// 空行のとき直前までのbufにstrを足したときに閾値を超えるか
					// 関数の末尾かっこにマッチさせるので、strも含める
					if IsOver(bufsize+len(str), num) {
						// 超えた
						buf = append(buf, str)
						out = append(out, buf)
						// リセット
						buf = []string{}
						bufsize = 0
					} else {
						// 超えていない、bufに結合する
						bufsize += len(str)
						buf = append(buf, str)
					}
				}
			}
		}
		// 残り
		out = append(out, buf)
	}
	return out
} // }}}

// サイズで区切る
// 空行のところでサイズチェックして超える前で分割する
// 空行がないとき、指定サイズを超える。
func DivideBySize(divided [][]string, num int) [][]string { // {{{
	out := [][]string{}

	for _, d := range divided {
		// dのサイズを計算
		s := calcStrArrayLength(d)
		//log.Printf("s:%v, num:%v", s, num)
		//log.Printf("-------------------------------------------------------------")
		if !IsOver(s, num) {
			//log.Printf("=============================================================")
			// 越えていないのでそのまま維持
			out = append(out, d)
			continue
		}
		//log.Printf("-------------------------------------------------------------")
		// サイズが超えているので分割
		flagPre := false // <pre>要素内で区切りたくない。
		// 一度bufに入れて、越えたときのbufをoutに吐く
		var buf []string // 越えていない状態の文字列群
		bufsize := 0
		for _, str := range d {
			// 行ごとの処理
			if reg_elem_pre.MatchString(str) { // <pre>に相当する```にマッチした
				flagPre = !flagPre
				log.Printf("reg_elem_pre.MatchString(str) にマッチ, str : %v, flagPre : %v", str, flagPre)
			}
			if flagPre { // <pre>要素の中の時
				bufsize += len(str)
				buf = append(buf, str)
			} else {
				// <pre>要素でない時
				// 空行があれば、その直前まででサイズを計算
				if !reg_elem_empty.MatchString(str) {
					bufsize += len(str)
					buf = append(buf, str) // 空行でないとき、bufに入れる。
				} else {
					// 空行のとき直前までのbufにstrを足したときに閾値を超えるか
					if IsOver(bufsize+len(str), num) {
						// 超えた
						out = append(out, buf)
						// リセット
						buf = []string{}
						bufsize = 0
					} else {
						// 超えていない、bufに結合する
						bufsize += len(str)
						buf = append(buf, str)
					}
				}
			}
		}
		// 残り
		out = append(out, buf)
	}
	return out
} // }}}

// heading要素 <h1> <h2> ... で区切る
func DivideOnHeadingElement(in [][]string) [][]string { // {{{
	// 区切り->行
	divided := [][]string{}
	buf := []string{}

	// とりあえず <h1,h2...> の要素で区切る。ただし、<pre>要素内で区切りたくない。
	flagPre := false
	for _, lines := range in {
		for _, line := range lines {
			if reg_elem_pre.MatchString(line) {
				flagPre = !flagPre
			}
			if !flagPre && // <pre>要素でない時
				reg_elem_h.MatchString(line) { // ^#にマッチしたら
				if len(buf) > 0 {
					//buf = append(buf, "\n") // 末尾に\nを追加
					divided = append(divided, buf)
					buf = []string{}
				}
			}
			buf = append(buf, line) // 末尾に追加
		}
		divided = append(divided, buf)
	}
	return divided
} // }}}

func DivideOnNewLine(text *string) [][]string { // {{{
	tmp := strings.Split(*text, "\n")
	// 改行しか入っていない要素は結合する。改行は消えているので長さ0になっている。
	//log.Printf("len(tmp) : %v", len(tmp))
	out := []string{}
	for _, t := range tmp {
		if len(t) != 0 {
			/*
				if len(t) > 10 {
					log.Printf("t : %v", t[len(t)-10:])
				} else {
					log.Printf("t : %v", t)
				}
			*/
			out = append(out, t)
		} else {
			if len(out) > 0 {
				out[len(out)-1] += "\n"
			} else {
				out = append(out, "\n")
			}
		}
	}
	//log.Printf("len(out) : %v", len(out))
	return [][]string{out}
} // }}}
