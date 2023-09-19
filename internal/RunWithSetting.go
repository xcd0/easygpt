package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/google/uuid"
)

type BoolString struct {
	status bool
	value  string
}
type BoolStrings []BoolString

func (bss *BoolStrings) Set(status bool, value *string) {
	for i, v := range *bss {
		if v.value == *value {
			mutex.Lock()
			(*bss)[i].status = status
			mutex.Unlock()
		}
	}
}

var mutex = &sync.Mutex{}
var failedFiles = BoolStrings{} // 処理に失敗したファイル

func RunWithSetting(setting *Setting) {
	fmt.Println("設定ファイルから実行しています")

	// 出力先ディレクトリを作成する。
	CreateOutputDirWithSetting(setting)

	// 入力ディレクトリの中のテキストファイルを検索
	// 拡張子指定があれば、その拡張子のファイルのみ使用する。
	files := GetTargetFiles(setting.InputDir, setting.Extension)
	if len(files) == 0 {
		fmt.Println("入力ディレクトリの中に処理対象となるテキストファイルがありませんでした。\n終了します。")
		os.Exit(0)
	}

	// 一時ファイルを保存するためのディレクトリを作成する。
	CreateTmpDir(setting)

	{
		if setting.Concurrency == 0 {
			setting.Concurrency = 1
		}
		fmt.Printf("最大平行処理数 : %v\n", setting.Concurrency)
	}

	// Askの準備
	// テキストを閾値で分割しておく。進捗バーに分割数を適用する。
	pbs := make([]*pb.ProgressBar, len(files))
	splited := map[string][]string{}
	for i := 0; i < len(files); i++ {
		f := &files[i]
		id := CreateID(f, setting)
		//log.Printf("debug id : %v", id)
		CreateDir(id)
		p := filepath.Join(id, "original_path.txt")
		OutputTextForCheck(&p, f)
		inputText := GetTextNoError(f)
		// テキストファイル分割
		strs := SplitText(inputText, setting.Split)
		log.Printf("split input text size : %v -> %v", f, len(strs))
		//log.Printf("split threshold : %v kB", setting.Split)
		//log.Printf("splited         : %v", len(strs))
		splited[*f] = strs
		// 2足しているのは、処理に入った瞬間に1を足している分と、結合処理の終了で1足している分
		//pbs[i] = pb.New(len(strs) + 3).Prefix(fmt.Sprintf("%-50s", filepath.Base(*f)))
		//pbs[i] = pb.New(len(strs) + 1).Prefix(fmt.Sprintf("%-50s", filepath.Base(*f)))
		pbs[i] = pb.New(len(strs)).Prefix(fmt.Sprintf("%-50s", filepath.Base(*f)))
	}

	log.Printf("split input text size : %v", len(splited))

	pool, err := pb.StartPool(pbs...)
	if err != nil {
		//panic(err) 無視
	}
	log.Printf("debug")

	// ファイルごとの処理
	// 並列数制限 semが解放されない限り、sem<-struct{}{}のところで処理が停止する。
	sem := make(chan struct{}, setting.Concurrency) // concurrency数のバッファ
	for i := 0; i < len(files); i++ {
		failedFiles = append(failedFiles, BoolString{false, files[i]})
	}
	var wg sync.WaitGroup
	for i := 0; i < len(files); i++ {
		sem <- struct{}{}
		wg.Add(1)

		f := files[i]
		p := pbs[i]
		//log.Printf("%v:%v", i, f)
		go func(f string, p *pb.ProgressBar) {
			defer func() {
				log.Printf("pool %v : %v finished", f, i)
				wg.Done()
				<-sem
				p.Finish()
			}()
			//defer func() { <-sem }()
			p.Increment()
			Ask(&f, setting, p, splited[f])
		}(f, p)
	}
	wg.Wait()
	pool.Stop()
	defer func() {
		if len(setting.Move) != 0 {
			RemoveEmptyDirectories(setting.InputDir) // 入力ディレクトリで空になっているディレクトリを削除
		}
	}()
}

func Move(src, dst string) {
	err := os.Rename(src, dst)
	if err != nil {
		log.Println("--")
		log.Printf("Error: Failed to move file: \nsrc: %v\ndst: -> %v", src, dst)
		log.Println(err)
	} else {
		log.Printf("File Moved: -> %v", dst)
	}
}

func Ask(input_file *string, setting *Setting, p *pb.ProgressBar, strs []string) bool {
	//defer p.Finish()
	//p.Increment()

	ret := false

	// 出力先ディレクトリパス
	outputpath := strings.Replace(*input_file, setting.InputDir, setting.OutputDir, 1)
	output := make([]string, len(strs))

	id := CreateID(input_file, setting)
	tmpflag := true
	if tmpflag {
		for i, str := range strs {
			//log.Printf("debug: inputText %d len : %v", i, len(str))
			// 分割したテキストを連番のファイル名で一時ディレクトリに保存する。
			p := filepath.Join(id, filepath.Base(*input_file)+"_"+strconv.Itoa(i)+".md")
			OutputTextForCheck(&p, &str)
		}
	}

	// 処理
	var wgAsk sync.WaitGroup
	for i, _ := range strs {
		wgAsk.Add(1)
		go func(i int, id *string, strs []string, setting *Setting, tmpflag bool, output *string) {
			outputText, err := task(i, id, &strs[i], setting, tmpflag)
			if err != nil {
				log.Printf("Error: %v, %+v", *input_file, err)
			} else {
				*output = *outputText + "\n"
				if tmpflag {
					// 分割したテキストを連番のファイル名で一時ディレクトリに保存する。
					p := filepath.Join(*id, filepath.Base(*input_file)+"_output_"+strconv.Itoa(i)+".md")
					OutputTextForCheck(&p, outputText)
				}
			}
			defer func() {
				p.Increment()
				wgAsk.Done()
				log.Printf("- go func defer() ----------------------------------\n\n")
			}()
		}(i, &id, strs, setting, tmpflag, &output[i])
	}
	wgAsk.Wait()

	log.Printf("\n-wgAsk.Wait()wgAsk.Wait()wgAsk.Wait()wgAsk.Wait()wgAsk.Wait()----------------------------------\n")
	log.Printf("\n-wgAsk.Wait()wgAsk.Wait()wgAsk.Wait()wgAsk.Wait()wgAsk.Wait()----------------------------------\n")

	if false {
		for i, o := range output {
			if len(o) > 15 {
				log.Printf("%v: '%s ... %s'", i, o[:15], o[len(o)-15:])
			} else {
				log.Printf("%v: '%s'", i, o)
			}
		}
	}

	// １回だけ再実行する。
	for i, str := range strs {
		if len(output[i]) == 0 {
			outputText, err := task(i, &id, &str, setting, tmpflag)
			if err == nil && len(output[i]) != 0 {
				output[i] = *outputText + "\n"
				if tmpflag {
					// 分割したテキストを連番のファイル名で一時ディレクトリに保存する。
					p := filepath.Join(id, filepath.Base(*input_file)+"_output_"+strconv.Itoa(i)+".md")
					OutputTextForCheck(&p, outputText)
				}
			} else {
				// このファイルは処理失敗とする。
				log.Printf("処理に失敗しました。 : %v", input_file)
				return ret
			}
		}
	}
	log.Printf("\n-結合処理----------------------------------\n")
	// 結合
	outputText := strings.Join(output, "\n")
	OutputTextForCheck(&outputpath, &outputText)
	log.Printf("\n-----------------------------------\n")

	// 正常終了。正常に処理完了したファイルを移動させる。
	ret = true
	defer func() {
		if ret && len(setting.Move) != 0 {
			if err := os.MkdirAll(setting.Move, 0777); err != nil {
				log.Println(err)
			}
			inrpath, _ := filepath.Rel(setting.InputDir, *input_file)
			outrpath, _ := filepath.Abs(filepath.Join(setting.Move, inrpath))
			if err := os.MkdirAll(filepath.Dir(outrpath), 0777); err != nil {
				log.Println(err)
			}
			Move(*input_file, outrpath)
			failedFiles.Set(true, input_file)
			log.Printf("\n-defer()----------------------------------\n")
		}
		defer func() {
			// 処理に失敗したファイルの一覧を出力する
			fmt.Println("処理に失敗したファイル:")
			for _, v := range failedFiles {
				if !v.status {
					fmt.Println(v.value)
				}
			}
			log.Printf("\n-defer()内defer()----------------------------------\n")
		}()
	}()
	return ret
}

func task(i int, id, str *string, setting *Setting, tmpflag bool) (*string, error) {
	d, err := os.MkdirTemp(*id, fmt.Sprintf("divided_%d_", i))
	if err != nil {
		log.Printf("%v", err)
	}
	return QuestionByText(str, setting, tmpflag, &d)
}

func RemoveEmptyDirectories(path string) error {
	emptyDirs := func(path string) ([]string, error) {
		var dirsToRemove []string
		if err := filepath.WalkDir(
			path,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				} else if d.IsDir() {
					if entries, err := os.ReadDir(path); err != nil {
						return err
					} else if len(entries) == 0 {
						dirsToRemove = append(dirsToRemove, path)
					}
				}
				return nil
			},
		); err != nil {
			return nil, err
		}
		return dirsToRemove, nil
	}
	dirsToRemove, err := emptyDirs(path)
	if err != nil {
		return err
	}
	for _, dir := range dirsToRemove {
		err := os.Remove(dir)
		if err != nil {
			fmt.Printf("Failed to remove directory: %v\n", dir)
		}
	}

	return nil
}
func CreateID(input_file *string, setting *Setting) string {
	// このファイルのファイルパスから固有の文字列を生成し、一時ディレクトリの名前とする。
	id := strings.ReplaceAll(filepath.ToSlash(strings.Replace(*input_file, setting.InputDir, "", 1)), "/", "-")[1:]
	// 一時ディレクトリにこのファイル固有のディレクトリを作成する
	id = filepath.Join(setting.Tmp, id)
	//log.Printf("debug")
	if len(id) == 0 {
		// なぜか文字がすべて消えた。ファイル名が\とかだったらありえる。
		// ランダム文字列生成してディレクトリ名にする
		uuidObj, _ := uuid.NewUUID()
		id = filepath.Join(setting.Tmp, uuidObj.String())
	}
	//log.Printf("debug id : %v", id)
	return id
}

func CreateOutputDirWithSetting(setting *Setting) {
	if errInputDir := CheckInputDir(&setting.InputDir); errInputDir != nil {
		fmt.Printf("%v\n\n", errInputDir)
		//parser.WriteUsage(os.Stdout)
		ShowUsage()
		os.Exit(1)
	}
	if errOutputDir := CheckOutputDir(&setting.OutputDir); errOutputDir != nil {
		fmt.Printf("%v\n\n", errOutputDir)
		//parser.WriteUsage(os.Stdout)
		ShowUsage()
		os.Exit(1)
	}

	// 出力先ディレクトリを作成
	CreateSameDirsOn(setting.InputDir, setting.OutputDir)
}

func CreateTmpDir(setting *Setting) {
	create_tmp := func() {
		d, err := os.MkdirTemp("", "eagygpt_tmp_")
		if err != nil {
			log.Printf("%v", err)
		} else {
			//log.Printf("tmp dir path : %v", d)
			setting.Tmp = fmt.Sprintf("%s", d)
		}
	}

	if len(setting.Tmp) == 0 {
		create_tmp()
	} else {
		os.RemoveAll(setting.Tmp)
		if err := os.Mkdir(setting.Tmp, 0777); err != nil {
			fe, e := "file exists", err.Error()
			if ee := e[len(e)-len(fe):]; ee != fe {
				log.Printf("Error: %v", err)
			}
			create_tmp() // 指定のディレクトリに一時ディレクトリを作成できなかった場合、OSの既定値の場所に作成する。
		}
	}
	ShowSetting(setting)
}
