package internal

import "fmt"

func (ArgsCommandLine) Description() string {
	return `# easygpt

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるアプリ。
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。
`
}

func ShowDescription() {
	fmt.Println(`# easygpt

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるアプリ。
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。

## 準備

設定ファイルの雛形を生成し、それを編集して設定ファイルを作成する。

1. 雛形の生成。

$ ./easygpt --create-setting

のように実行するとカレントディレクトリに eagygpt.hjson が生成される。  

2. 設定ファイルの編集。

* "apikey" の部分にAPIキーを書き込む。  
	例) "apikey: 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"  
	これは無効なAPIキー。自分で発行する事。  
	APIキーは https://platform.openai.com/account/api-keys から発行できる。

* "prompt" の部分に入力テキストファイルの前に与えたい文字列を記載する。  
	例) 英文テキストファイルを翻訳してほしい場合、"prompt: 以下を和訳してください。"とのように書く。  
	複数行書きたいとき、"\n"で改行するか、以下のヒアドキュメント形式で記載する。  
	prompt:  
	'''  
	このような書式で、  
	複数行書くことができる。  
	'''  

他は説明コメントを読みながら必要に合わせて記述する。

3. 設定ファイルの配置。

以下の仕様に従って設定ファイルが探索される。  
基本的にはカレントディレクトリに、easygpt.hjsonを配置すればよい。

設定ファイルの名前の既定値は、  
* easygpt.hjson  
* .easygpt.hjson  
* .easygpt  
の3種類。  
これ以外の名前の設定ファイルを使用したい場合、引数"--setting"で設定できる。  

設定の既定の位置は、  
* カレントディレクトリ  
* ホームディレクトリ  
* 実行ファイルと同じディレクトリ  
の3か所。  
設定ファイルの探索は、この順番で行われる。  
従って、カレントディレクトリに存在すれば、ホームディレクトリの設定ファイルは無視される。  

## 使い方1 D&D

設定ファイルがあれば、以下のように使用できる。

1. gptに投げたいテキストファイル、またはそれが含まれるディレクトリを、 "easygpt"の実行ファイルにドラッグアンドドロップする。
1. 投げたファイルと同じディレクトリに、入力ファイルに"_easygpt_output"を付与した名前で処理結果を出力される。

## 使い方2 コマンドラインから実行

設定ファイルの設定をコマンドラインから指定できる。  
既定のパスにある設定ファイルに、設定されている値については、引数を省略できる。  
コマンドライン引数が設定されていた場合、設定ファイルの値をコマンドライン引数の値で上書きする。  
設定できる引数が多い為、基本的に設定ファイルを使用する事をお勧めする。  

引数には "--key value" 形式のフラグ引数と、そうでない非フラグ引数がある。  
* フラグ引数は、優先度があり、特定の引数を使用した場合、他の引数を無視する場合がある。  
  例えば、"--create-setting" があれば、他の引数を無視して、設定ファイルの雛形を生成した後終了する。  
  value を空にしたい場合を除いて、""で囲む必要はない。  
* 非フラグ引数は、入力ファイルと見なして処理される。  
  ディレクトリを指定した場合、そのディレクトリを再帰的に探索して全てのファイルを入力ファイルとする。  
引数は順不同。  
全ての引数を設定すると以下の例のようになる。  

個々の引数の詳細については "./easygpt --help" を参照。  

$ ./easygpt                                                                \  
                     "./input_1.txt" "./input_2.txt"                       \  
                     "./input_dir_1" "./inputdir_2"                        \  
    --create-setting "./easygpt.hjson"                                     \  
    --setting        "./easygpt.hjson"                                     \  
    --input-dir      "./inputdir"                                          \  
    --output-dir     "./outputdir"                                         \  
    --apikey-file    "./apikey.txt"                                        \  
    --apikey         "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" \  
    --prompt-file    "./prompt.txt"                                        \  
    --prompt         "以下を和訳してください。"                              \  
    --postfix        "_easygpt_output"                                     \  
    --extension      ""                                                    \  
    --tmp-dir        "./tmpdir"                                            \  
    --concurrency    1  

### 例

./inputに英文テキストファイルがあるとして、それらをまとめて和訳させる例。
./outputに出力が保存される。

$./easygpt \
	--input-dir ./input \
	--output-dir ./output \
	--api-key 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 \
	--prompt 以下を和訳してください。

例えば、設定ファイル"./easygpt.hjson"にAPIキーが設定されていた場合、上記の引数"--api-key"は省略できる。  

`)
}
