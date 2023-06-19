# easygpt

[English readme translated by this program.](./readme_en.md)

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるプログラム。  
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。  

## install

```sh
go install github.com/xcd0/easygpt@latest
```

## 準備

OpenAIのAPIキーを環境変数`OPENAI_API_KEY`に設定する。  
下記のAPIキーは例の為の無効なAPIキーなので、自分で有効なAPIキーを発行して、設定すること。  
APIキーは https://platform.openai.com/account/api-keys から発行できる。
発行したAPIキーを`.bashrc`などに設定する。
```sh
$ export OPENAI_API_KEY=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
$ echo export OPENAI_API_KEY=$OPENAI_API_KEY >> .bashrc
```

## 使用例1 : AIに文字列を与えて出力を得る

APIキーを環境変数で指定していれば、コマンドラインから直接実行できる。

```sh
$ easygpt -i "Say this is test."
this is test.
$ easygpt -i 自己紹介してください。
こんにちは、私はAIです。私はOpenAIが開発した自然言語処理モデルです。私の目的は、ユーザーが質問や要求をすると、最善の回答や応答を提供することです。私は様々なトピックについての情報を持っており、文法やスタイルの修正も行うことができます。どのようにお手伝いできますか？
```

## 使用例2 : D&D

与えたファイルとフォルダに含まれるファイル全てを一括でAIに与えて処理させる。  
ファイルやフォルダを選択して実行ファイルにD&Dすることで、含まれるファイルをテキストとしてAIに処理させることができる。  
* GUI上での使用を想定。  
* 出力は入力ファイルに`_easygpt_output`を付与した名前で出力される。  
* 後述の設定ファイルによって動作を変更できる。
	* この使い方では設定ファイルの設定のうち、
		* input-dir
		* output-dir
		* extension
		を無視する。

## 使用例3 : 複数ファイル一括処理

AIに指定のディレクトリにある複数のテキストファイルを一括で処理させる。  
これがメインの使い方。  
この使い方は、すべて設定ファイルで動作を指定する。  
使用例2よりも細かい制御ができる。  

## 設定ファイルについて

設定ファイルの雛形を生成し、それを編集して設定ファイルを作成する。  

| key          | 必須か | 説明                  |
|--------------|--------|-----------------------|
| input-dir    | 必須   | 入力ディレクトリ。<br>AIに与えたいファイル群を格納するディレクトリ。 |
| output-dir   | 必須   | 出力先ディレクトリ。<br>入力ディレクトリと同じディレクトリ構成になるようにディレクトリが作成される。|
| apikey       | 任意   | APIキーの指定。環境変数に指定していれば不要。<br>設定ファイルで指定されていれば環境変数での指定を無視する。|
| prompt       | 任意   | ここで指定した文字列を、全ての入力ファイルの先頭に付与する。<br>AIに対する支持などを記載する。|
| extension    | 任意   | 入力ディレクトリにおいて指定の拡張子で入力ファイルを制限する。<br> この指定が空の時、*を指定した時、無指定の時、拡張子によって制限しない。|
| concurrency  | 任意   | 平行して実行する数。APIのレート制限に注意。|
| temperature  | 任意   | AIの変数`temperature`を指定できる。 |
| move         | 任意   | 正常に処理ができたファイルを指定したディレクトリに移動させる。<br>エラーやCtrl-Cなどで中断した際に、途中から再度実行しやすくなる。 |
| tmp-dir      | 任意   | 基本的に指定する必要はない。一時ファイルを保持するディレクトリ。|
| postfix      | 任意   | 基本的に指定する必要はない。出力ファイル名の末尾に文字列を付与したいとき指定する。|
| ai-model     | 任意   | 基本的に指定する必要はない。使用したいAIのモデルを指定できる。|
| openai-url   | 任意   | 基本的に指定する必要はない。APIのURLを指定できる。|


1. 雛形の生成。

設定できる項目が多い為、設定ファイルの雛形を生成し、それを編集することを推奨。  

以下のコマンドで、カレントディレクトリに `eagygpt.hjson` が生成できる。  
```sh
$ easygpt --create-setting
```
`easygpt -c`でも良い。


2. 設定ファイルの編集。
* APIキーの設定
	* 設定ファイルにAPIキーを指定できる。環境変数`OPENAI_API_KEY`に設定していれば不要。

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
	詳しくは[hjsonのドキュメント](https://hjson.github.io/syntax.html)の`Mulutiline Setting`を参照。

他の引数については、生成された設定ファイルのコメントを読み、必要に合わせて記述する。  
基本的には、`input-dir`と`output-dir`が指定されていればよい。

3. 設定ファイルの配置。

基本的にはカレントディレクトリに、easygpt.hjsonを配置すればよい。  

以下の仕様に従って設定ファイルが探索される。  

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


## 引数による設定

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
		* makeの際にupxで圧縮し、無改変で配布している。プログラムにupxのソースコードは含まない。

## LICENSE

MIT License


