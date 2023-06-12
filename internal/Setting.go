package internal

type Setting struct {
	Apikey      string `json:"apikey"      comment:"APIキーを指定する。"`                                                                                                                                                                                                                      // APIキー
	InputDir    string `json:"input-dir"   comment:"入力テキストファイルのあるディレクトリパス。\n再帰的にファイルを検索して全て処理する。"`                                                                                                                                                                                     // 入力ディレクトリ
	OutputDir   string `json:"output-dir"  comment:"出力テキストファイルのあるディレクトリパス。\ninput-dirのディレクトリと同じ構造でサブディレクトリを作成する。"`                                                                                                                                                                     // 出力ディレクトリ
	Prompt      string `json:"prompt"      comment:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\nこの指定があるときファイルによる指定を無視する。\nプロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n例えば、「以下のテキストを和訳してください。」と書き、入力に英文テキストファイルを与えれば翻訳して貰える。\nヒアドキュメントのように複数を見やすく記述することもできる。\n'''\nヒアドキュメント形式の\n複数行プロンプト\n記述例\n'''"` // プロンプト
	Extension   string `json:"extension"   comment:"入力として使用したいテキストファイルの拡張子。\n拡張子のドットを含めて.mdのように指定する。指定なしの時すべて使用する。"`                                                                                                                                                                   // 拡張子
	Concurrency int    `json:"concurrency" comment:"並列処理数を指定する。初期値1。\nAPIのRateLimitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\nこの数値は単純に並行処理のスレッド数だと思ったらよい。Token/分とRequest/分に配慮する。\n小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。"`                                                              // 最大同時並列実行数
	Tmp         string `json:"tmp-dir"     comment:"一時ファイルを保存するディレクトリを指定する。\n指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n既にあれば、削除して再作成する。"`                                                                                                                                                     // 一時ファイルを保存するディレクトリ
	Postfix     string `json:"postfix"     comment:"出力ファイル名の末尾に付与する文字列。"`                                                                                                                                                                                                              // 出力ファイル名の末尾に付与する文字列。空の時 "_easygpt_output" となる。
}

type SettingForDD struct {
	ApiKey      string // APIキー
	Input       string // 入力ファイル
	Prompt      string // プロンプト
	Postfix     string // 出力ファイル名の末尾に付与する文字列。空の時 "_easygpt_output" となる。
	Tmp         string // 一時ファイルを保存するディレクトリ
	Concurrency int    // 最大同時並列実行数
}
