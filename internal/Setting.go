package internal

type Setting struct {
	ApiKey      string // APIキー
	InputDir    string // 入力ディレクトリ
	OutputDir   string // 出力ディレクトリ
	Prompt      string // プロンプト
	Postfix     string // 出力ファイル名の末尾に付与する文字列。空の時 "_easygpt_output" となる。
	Extension   string // 拡張子
	Tmp         string // 一時ファイルを保存するディレクトリ
	Concurrency int    // 最大同時並列実行数
}

type SettingForDD struct {
	ApiKey      string // APIキー
	Input       string // 入力ファイル
	Prompt      string // プロンプト
	Postfix     string // 出力ファイル名の末尾に付与する文字列。空の時 "_easygpt_output" となる。
	Tmp         string // 一時ファイルを保存するディレクトリ
	Concurrency int    // 最大同時並列実行数
}
