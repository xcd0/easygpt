package internal

type ArgsAll struct {
	ArgsSetting
	ArgsDD
	ArgsCommandLine

	InputText  string `arg:"-i,--input-text"  help:"質問したい文字列を直接指定。\n                         この引数がある場合、APIキー指定、OpenAIのURL指定、AIのモデル指定、AIのTemperature、\n                         これら以外の引数と設定を無視します。\n"`
	Readme     bool   `arg:"-r,--readme"      help:"詳しい説明文を出力する。長いため\"./easygpt --readme | less\"などで見るのがおすすめ。\n                         内容はreadme.mdと同じ。"`
	OutputText string `arg:"-o,--output-text" help:"--input-textの出力を指定のテキストファイルに保存する。\n                         --input-textの指定がない時この指定は無視される。\n"`
}

type ArgsSetting struct {
	CreateSetting bool   `arg:"-c,--create-setting"   help:"設定ファイルの雛形を生成する。\n"`
	Setting       string `arg:"-s,--setting"          help:"設定ファイルを指定する。\n                         指定がなければカレントディレクトリか、ホームディレクトリか、実行ファイルのあるディレクトリを探す。\n                         探すファイル名は'easygpt.hjson'、'.easygpt.hjson'、'.easygpt'の3つ。\n"`
}

type ArgsDD struct {
	InputFiles []string `arg:"positional"         help:"非フラグ引数はファイルやディレクトリのパスであると見なす。\n                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、\n                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。"`
}

type ArgsCommandLine struct {
	InputDir    string  `arg:"--input-dir"       help:"入力テキストファイルのあるディレクトリパス。\n                         再帰的にファイルを検索して全て処理する。\n"`
	OutputDir   string  `arg:"--output-dir"      help:"出力テキストファイルのあるディレクトリパス。\n                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。\n"`
	ApiFile     string  `arg:"--apikey-file"     help:"テキストファイルによってAPIキーを指定する。\n                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。\n"`
	Apikey      string  `arg:"-a,--apikey"       help:"APIキーを直接指定する。\n                         この指定がある時、--api-fileによる指定を無視する。\n"`
	PromptFile  string  `arg:"--prompt-file"     help:"全ての入力テキストファイルの直前に付与したい文字列を書いたテキストファイルパス。\n                         例えば、「以下のテキストを和訳してください。」と書いたテキストファイルを指定して、\n                         入力テキストファイルとして英文テキストファイルを与えれば、翻訳して貰える。\n                         この指定がない時、カレントディレクトリの./prompt.txtがあればそれを使用する。\n"`
	Prompt      string  `arg:"-p,--prompt"       help:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\n                         この指定があるときファイルによる指定を無視する。\n                         プロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n"`
	Postfix     string  `arg:"--postfix"         help:"出力ファイル名の末尾に付与する文字列。\n                         空の時 _easygpt_output となる。\n"`
	Extension   string  `arg:"-e,--extension"    help:"入力として使用したいテキストファイルの拡張子。\n                         指定なしの時すべて使用する。\n                         拡張子のドットを含めて.mdのように指定する。\n"`
	Tmp         string  `arg:"--tmp-dir"         help:"一時ファイルを保存するディレクトリを指定する。\n                         指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n                         既にあれば、削除して再作成する。\n"`
	Concurrency int     `arg:"--concurrency"     help:"並列処理数を指定する。初期値1。\n                         APIの Rate Limitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\n                         この数値は単純に並行処理のスレッド数だと思ったらよい。\n                         Token/分とRequest/分に配慮する。\n                         小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。\n"`
	AiModel     string  `arg:"-m,--model"        help:"使用するAIのモデル。\n                         3か月くらいで新しいモデルが出るので偶にチェックするのがおすすめ。\n                         使用できるモデルは\n                         curl -s https://api.openai.com/v1/models -H \"Authorization: Bearer $OPENAI_API_KEY\" | gojq -r \".data[].id\" | grep gpt | sort\n                         で得られる。\n"`
	Temperature float64 `arg:"-t,--temperature"  help:"これは、AIに与える変数で、返答のランダム性を制御するパラメータである。\n                         値が小さいほどよくある解答など決まりきった解答を返し、\n                         値が大きいほど奇抜な返答が返ってきやすくなる。\n                         0から2の範囲の値を設定する。詳細は\n                         https://platform.openai.com/docs/api-reference/completions/create#completions/create-temperature\n                         を参照。\n"`
	OpenaiURL   string  `arg:"--openai-url"      help:"OpenAIのAPIで使用するURL。\n                         基本的に\"https://api.openai.com/v1/chat/completions\"だが、\n                         将来変更されるかもしれないので設定できるようにしておく。\n"`
	Split       int     `arg:"--split"           help:"入力ファイルが大きいため、AIの出力が途切れてしまう場合がある。\n                         これを回避するために、入力ファイルを分割してAIに与えたいとき、\n                         ここに分割するバイト数をkB単位で大雑把に指定できる。\n                         0~16の値を指定できる。単位はkB。トークン数ではない。0または未指定の時、分割しない。\n"`
	Move        string  `arg:"-m,--move"         help:"処理が完了したファイルを指定のディレクトリに移動する。\n                         エラー時に再実行しやすくなる。未指定の時移動させない。\n"`
}
