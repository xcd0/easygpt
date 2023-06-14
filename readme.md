# easygpt

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるアプリ。  
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。  

## install

```sh
go install github.com/xcd0/easygpt@latest
```

## 準備

設定ファイルの雛形を生成し、それを編集して設定ファイルを作成する。  
設定ファイル無しでも実行可能ではあるが、

1. 雛形の生成。

```sh
$ ./easygpt --create-setting
```

のように実行するとカレントディレクトリに `eagygpt.hjson` が生成される。  

2. 設定ファイルの編集。

* `apikey` の部分にAPIキーを書き込む。  
	例) `apikey: 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`  
	上記は例の為の無効なAPIキーであるので、自分で発行する事。  
	APIキーは https://platform.openai.com/account/api-keys から発行できる。

* `prompt` の部分に入力テキストファイルの前に与えたい文字列を記載する。  
	例) 英文テキストファイルを翻訳してほしい場合、`prompt: 以下を和訳してください。` のように書く。  
	複数行書きたいとき、`\n`で改行するか、以下のヒアドキュメント形式で記載する。  
	```
	prompt:  
	'''  
	このような書式で、  
	複数行書くことができる。  
	'''  
	```

他の引数については、設定ファイルのコメントを読み、必要に合わせて記述する。  

3. 設定ファイルの配置。

以下の仕様に従って設定ファイルが探索される。  
基本的にはカレントディレクトリに、easygpt.hjsonを配置すればよい。

設定ファイルの名前の既定値は、  
* easygpt.hjson  
* .easygpt.hjson  
* .easygpt  
の3種類。  
これ以外の名前の設定ファイルを使用したい場合、引数`--setting`で設定できる。  

設定の既定の位置は、  
* カレントディレクトリ  
* ホームディレクトリ  
* 実行ファイルと同じディレクトリ  
の3か所。  
設定ファイルの探索は、この順番で行われる。  
従って、カレントディレクトリに存在すれば、ホームディレクトリの設定ファイルは無視される。  

## 使い方1 D&D

設定ファイルがあれば、以下のように使用できる。  

1. gptに投げたいテキストファイル、またはそれが含まれるディレクトリを、 `easygpt`の実行ファイルにドラッグアンドドロップす
る。
1. 投げたファイルと同じディレクトリに、入力ファイルに`_easygpt_output`を付与した名前で処理結果を出力される。

## 使い方2 コマンドラインから実行

設定ファイルの設定をコマンドラインから指定できる。  
既定のパスにある設定ファイルに、設定されている値については、引数を省略できる。  
コマンドライン引数が設定されていた場合、設定ファイルの値をコマンドライン引数の値で上書きする。  
設定できる引数が多い為、基本的に設定ファイルを使用する事をお勧めする。  

引数には `--key value` 形式のフラグ引数と、そうでない非フラグ引数がある。  
* フラグ引数は、優先度があり、特定の引数を使用した場合、他の引数を無視する場合がある。  
  例えば、`--create-setting` があれば、他の引数を無視して、設定ファイルの雛形を生成した後終了する。  
  value を空にしたい場合を除いて、""で囲む必要はない。  
* 非フラグ引数は、入力ファイルと見なして処理される。  
  ディレクトリを指定した場合、そのディレクトリを再帰的に探索して全てのファイルを入力ファイルとする。  

引数は順不同。  
全ての引数を設定すると以下の例のようになる。  
`#` の後ろはコメント  
これはあくまで例にすぎず、実際には`--help`が優先されてヘルプテキストが出力される。  

```sh
$ ./easygpt                                                                \
                     "./input_1.txt" "./input_2.txt"                       \  # 非フラグ引数で入力ファイル指定
                     "./input_dir_1" "./inputdir_2"                        \  # 非フラグ引数で入力ファイルのあるディレクトリ指定
    --create-setting "./easygpt.hjson"                                     \  # 設定ファイルの雛形を生成
    --setting        "./easygpt.hjson"                                     \  # 読み込む設定ファイルのパスを指定
    --input-dir      "./inputdir"                                          \  # 入力ディレクトリ指定
    --output-dir     "./outputdir"                                         \  # 出力ディレクトリ指定
    --apikey-file    "./apikey.txt"                                        \  # APIキーを書いたテキストファイル指定
    --apikey         "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" \  # APIキーを直接指定
    --prompt-file    "./prompt.txt"                                        \  # プロンプトを書いたファイル指定
    --prompt         "以下を和訳してください。"                              \  # プロンプトを直接指定
    --postfix        "_easygpt_output"                                     \  # 出力ファイルのファイル名の末尾に付与する文字列指定
    --extension      ""                                                    \  # 入力ファイルを指定の拡張子に制限する
    --tmp-dir        "./tmpdir"                                            \  # 処理で使用する一時ディレクトリの出力先指定
    --concurrency    1                                                     \  # 並行処理数の指定
    --readme                                                                  # 詳しい説明文の出力
    --help                                                                    # ヘルプテキストの出力
```

引数の詳細については `./easygpt -h` で出力されるヘルプテキストを参照。  

以下は`./input`に英文テキストファイルがあるとして、それらをまとめて和訳させる例。

```
$./easygpt \
    --input-dir ./input \
    --output-dir ./output \
    --api-key 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 \
    --prompt 以下を和訳してください。
```

例えば、設定ファイル`./easygpt.hjson`にAPIキーが設定されていた場合、上記の引数`--api-key`は省略できる。  

## ヘルプ

```sh
$ ./easygpt --help
# easygpt

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるアプリ。
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。

Usage: easygpt [--create-setting CREATE-SETTING] [--setting SETTING] [--input-dir INPUT-DIR] [--output-dir OUTPUT-DIR] [--apikey-file APIKEY-FILE] [--apikey APIKEY] [--prompt-file PROMPT-FILE] [--prompt PROMPT] [--postfix POSTFIX] [--extension EXTENSION] [--tmp-dir TMP-DIR] [--concurrency CONCURRENCY] [--readme] [INPUTFILES [INPUTFILES ...]]

Positional arguments:
  INPUTFILES             非フラグ引数はファイルやディレクトリのパスであると見なす。
                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、
                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。

Options:
  --create-setting CREATE-SETTING
                         設定ファイルの雛形を生成する。

  --setting SETTING      設定ファイルを指定する。
                         指定がなければカレントディレクトリか、ホームディレクトリか、実行ファイルのあるディレクトリを探す。
                         探すファイル名は'easygpt.hjson'、'.easygpt.hjson'、'.easygpt'の3つ。

  --input-dir INPUT-DIR
                         入力テキストファイルのあるディレクトリパス。
                         再帰的にファイルを検索して全て処理する。

  --output-dir OUTPUT-DIR
                         出力テキストファイルのあるディレクトリパス。
                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。

  --apikey-file APIKEY-FILE
                         テキストファイルによってAPIキーを指定する。
                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。

  --apikey APIKEY        APIキーを直接指定する。
                         この指定がある時、--api-fileによる指定を無視する。

  --prompt-file PROMPT-FILE
                         全ての入力テキストファイルの直前に付与したい文字列を書いたテキストファイルパス。
                         例えば、「以下のテキストを和訳してください。」と書いたテキストファイルを指定して、
                         入力テキストファイルとして英文テキストファイルを与えれば、翻訳して貰える。
                         この指定がない時、カレントディレクトリの./prompt.txtがあればそれを使用する。

  --prompt PROMPT        全ての入力テキストファイルの直前に付与したい文字列を指定する。
                         この指定があるときファイルによる指定を無視する。
                         プロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。

  --postfix POSTFIX      出力ファイル名の末尾に付与する文字列。
                         空の時 _easygpt_output となる。

  --extension EXTENSION
                         入力として使用したいテキストファイルの拡張子。
                         指定なしの時すべて使用する。
                         拡張子のドットを含めて.mdのように指定する。

  --tmp-dir TMP-DIR      一時ファイルを保存するディレクトリを指定する。
                         指定がない時、カレントディレクトリにtmpディレクトリを作成する。
                         既にあれば、削除して再作成する。

  --concurrency CONCURRENCY
                         並列処理数を指定する。初期値1。
                         APIの Rate Limitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。
                         この数値は単純に並行処理のスレッド数だと思ったらよい。
                         Token/分とRequest/分に配慮する。
                         小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。

  --readme               詳しい説明文を出力する。長いため"./easygpt --readme | less"などで見るのがおすすめ。
  --help, -h             display this help and exit
```


## 使用させていただいているOSSライブラリ

各OSSライブラリのライセンスは `./lisenses` にコピーを保持している。

* 引数解析、ヘルプ生成
    * [alexflint/go-arg](https://github.com/alexflint/go-arg)
        * BSD-2-Clause license
* ランダム文字列生成 (別にUUIDである必要はない)
    * [google/uuid](https://github.com/google/uuid)
        * BSD-3-Clause license
* スタックトレース
    * [pkg/errors](https://github.com/pkg/errors)
        * BSD 2-Clause "Simplified" License
* golang本体
    * [golang/go](https://github.com/golang/go)
        * BSD 3-Clause "New" or "Revised" License
* ライセンスファイル群取得
    * [google/go-licenses](https://github.com/google/go-licenses)
        * Apache License 2.0
* 実行ファイル圧縮
    * [upx/upx](https://github.com/upx/upx)
        * GPL2+ or UPX LICENSE
            * https://upx.github.io/upx-license.html
        * makeの際にバイナリを使用している。かつ使用後のバイナリを無改変で配布している。プログラムにupxのソースコードは含まない。

## LICENSE

MIT License


