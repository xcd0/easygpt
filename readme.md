# easygpt

[English translated by chatgpt](./readme_en.md)

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるアプリ。  
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。  

## install

```sh
go install github.com/xcd0/easygpt@latest
```

## 単純な使用例 (AIに文字列を与えて出力を得る)

OpenAIのAPIキーを環境変数`OPENAI_API_KEY`に設定する。  
下記のAPIキーは例の為の無効なAPIキーなので、自分で有効なAPIキーを発行して、設定すること。  
APIキーは https://platform.openai.com/account/api-keys から発行できる。
発行したAPIキーを`.bashrc`などに設定する。
```sh
$ export OPENAI_API_KEY=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
$ echo export OPENAI_API_KEY=$OPENAI_API_KEY >> .bashrc
```

あとはeasygptをインストールして、実行すればよい。

```sh
$ go install github.com/xcd0/easygpt@latest
$ easygpt -i "Say this is test."
this is test.
```

## 複数ファイル一括処理

AIに複数のテキストファイルを一括で処理させる。  

### 準備

設定ファイルの雛形を生成し、それを編集して設定ファイルを作成する。  

1. 雛形の生成。

```sh
$ easygpt --create-setting
```

のように実行すると、カレントディレクトリに `eagygpt.hjson` が生成される。  
`easygpt -c`でも良い。

2. 設定ファイルの編集。
* APIキーの設定
	* 環境変数`OPENAI_API_KEY`に設定しておくのがおすすめ。この場合設定ファイルに書き込む必要はない。
	* 設定ファイルの`apikey` の部分にAPIキーを書き込んでもよい。  
		例) `apikey: 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`  
		上記のAPIキーは無効なAPIキーであるので、自分で発行する事。  
		APIキーは https://platform.openai.com/account/api-keys から発行できる。

* 設定ファイルの`prompt` の部分に、入力テキストファイルの前に与えたい文字列を記載する。  
	例) 英文テキストファイルを翻訳してほしい場合、
	`prompt: 以下を和訳してください。` のように書く。  
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

GUI上での使用方法。ファイルやフォルダを選択して実行ファイルにD&Dする。  
与えたファイルとフォルダに含まれるファイル全てを一括でAIに与えて処理させる。  
出力は入力ファイルに`_easygpt_output`を付与した名前で出力される。  

設定ファイルがあれば、以下のように使用できる。  
1. gptに投げたいテキストファイル、またはそれが含まれるディレクトリを、 
`easygpt`の実行ファイルにD&Dする。
1. 入力ファイルに`_easygpt_output`を付与した名前で処理結果を出力される。  

この使い方の時は設定ファイルの設定のうち、
* input-dir
* output-dir
* extension
を無視する。

## 使い方2 コマンドラインから実行

簡単な実行方法の例として、環境変数や、設定ファイルにAPIキーが設定されていれば、以下のように使用できる。
```sh
$ ./easygpt -i 自己紹介してください。
こんにちは、私はAIです。私はOpenAIが開発した自然言語処理モデルです。私の目的は、ユーザーが質問や要求をすると、最善の回答や応答を提供することです。私は様々なトピックについての情報を持っており、文法やスタイルの修正も行うことができます。どのようにお手伝いできますか？
```

## 使い方3 コマンドラインから実行 ファイル群を処理させる

コマンドライン引数が設定されていた場合、設定ファイルの値をコマンドライン引数の値で上書きする。  
既定のパスにある設定ファイルに、設定されている値については、引数を省略できる。  
設定できる引数が多い為、基本的に設定ファイルを使用する事をお勧めする。  

引数には `--key value` 形式のフラグ引数と、そうでない非フラグ引数がある。  
* フラグ引数は、優先度があり、特定の引数を使用した場合、他の引数を無視する場合がある。  
  例えば、`--create-setting` があれば、他の引数を無視して、設定ファイルの雛形を生成した後終了する。  
  value を空にしたい場合を除いて、""で囲む必要はない。  
* 非フラグ引数は、入力ファイルと見なして処理される。  
  ディレクトリを指定した場合、そのディレクトリを再帰的に探索して全てのファイルを入力ファイルとする。  

引数は順不同。  
引数の詳細については `./easygpt --help` で出力されるヘルプテキストを参照。  

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
* 設定ファイル
	* [hjson/hjson-go](https://github.com/hjson/hjson-go)
		* MIT License
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


