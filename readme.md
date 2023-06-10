# easygpt

テキストファイルをchatgptのapiに投げる。  
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。

## install

```sh
go install github.com/xcd0/easygpt@latest
```

## 使い方1 D&D

1. 実行ファイルと同じディレクトリに2つテキストファイルを作成する。
	* `./apikey.txt`
		* APIキーを書き込む。  
		例) `echo "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" > ./apikey.txt`
            * このAPIキーは無効なのでちゃんと自分のアカウントで発行して設定すること。
		* APIキーは https://platform.openai.com/account/api-keys から発行できる。
	* `./prompt.txt`
		* これはなくてもよい。
		* 入力テキストファイルの前に与えたい文字列を記載する。
			* 例) 英文テキストファイルを翻訳してほしい場合、`./prompt.txt`に`以下を和訳してください。`と書く。

1. gptに投げたいテキストファイル、またはそれが含まれるディレクトリを`easygpt.exe`の実行ファイルにドラッグアンドドロップする。
1. 投げたファイルと同じディレクトリに、入力ファイルに`_easygpt_output`を付与した名前で処理結果を出力する。


## 使い方2 コマンドラインから実行

詳しくは`./easygpt -h`を参照

1. ヘルプを見ながら引数を設定する。

### 例

./inputに英文テキストファイルがあるとして、それらをまとめて和訳させる例。

```sh
$ ./easygpt \
	--input-dir "./input" \
	--output-dir "./output" \
	--api-key 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 \
	--prompt "以下を和訳してください。" 
```

## help

```
$ ./easygpt -h
this program does this and that
Usage: easygpt [--input-dir INPUT-DIR] [--output-dir OUTPUT-DIR] [--api-file API-FILE] [--api-key API-KEY] [--prompt-file PROMPT-FILE] [--pr
ompt PROMPT] [--postfix POSTFIX] [--target-extension TARGET-EXTENSION] [--tmp-dir TMP-DIR] [--concurrency CONCURRENCY] [INPUTFILES [INPUTFIL
ES ...]]

Positional arguments:
  INPUTFILES             非フラグ引数はファイルやディレクトリのパスであると見なす。
                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、
                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。

Options:
  --input-dir INPUT-DIR
                         入力テキストファイルのあるディレクトリパス。
                         再帰的にファイルを検索して全て処理する。

  --output-dir OUTPUT-DIR
                         出力テキストファイルのあるディレクトリパス。
                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。

  --api-file API-FILE    テキストファイルによってAPIキーを指定する。
                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。

  --api-key API-KEY      APIキーを直接指定する。
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

  --target-extension TARGET-EXTENSION
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


