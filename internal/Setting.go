package internal

type Setting struct {
	InputDir  string `json:"input-dir"   comment:"入力テキストファイルのあるディレクトリパス。\n再帰的にファイルを検索して全て処理する。\n例) input-dir: ./input_dir"`
	OutputDir string `json:"output-dir"  comment:"出力テキストファイルのあるディレクトリパス。\ninput-dirのディレクトリと同じ構造でサブディレクトリを作成する。\n例) output-dir: ./output_dir"`
	SettingCommon
}

type SettingForDD struct {
	Input     string // 入力ファイル
	OutputDir string

	Common *SettingCommon
}

type SettingCommon struct {
	Apikey      string  `json:"apikey"      comment:"APIキーを指定する。\nAPIキーは https://platform.openai.com/account/api-keys から発行できる。\n例) apikey: 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"`
	Prompt      string  `json:"prompt"      comment:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\nこの指定があるときファイルによる指定を無視する。\nプロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n例えば、「以下のテキストを和訳してください。」と書き、入力に英文テキストファイルを与えれば翻訳して貰える。\nヒアドキュメントのように複数を見やすく記述することもできる。\nprompt:\n\t'''\n\tヒアドキュメント形式の\n\t複数行プロンプト\n\t記述例\n\t'''\n例) prompt: \"以下を翻訳してください。\"\n(必要でなければ設定不要)"`
	Extension   string  `json:"extension"   comment:"入力として使用したいテキストファイルの拡張子。\n拡張子のドットを含めて.mdのように指定する。拡張子に関わりなく全て使用したい場合 \"\" のように空にする。\n例) extension: .md # markdownのファイルだけを処理させたい場合の例\n例) extension: \"\"    # 拡張子で入力ファイルを制限しない場合の例\n(必要でなければ設定不要)"`
	Concurrency int     `json:"concurrency" comment:"並列処理数を指定する。初期値1。\nAPIのRateLimitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\nこの数値は単純に並行処理のスレッド数だと思ったらよい。Token/分とRequest/分に配慮する。\n小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。\n例) concurrency: 4\n(必要でなければ設定不要)"`
	Tmp         string  `json:"tmp-dir"     comment:"一時ファイルを保存するディレクトリを指定する。\n指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n既にあれば、削除して再作成する。\n例) tmp-dir: ./tmp_dir\n(基本的に設定不要)"`
	Postfix     string  `json:"postfix"     comment:"出力ファイル名の末尾に付与する文字列。\n例) postfix: _easygpt_output\n(必要でなければ設定不要)"`
	AiModel     string  `json:"ai-model"    comment:"使用するAIのモデル。\n3か月くらいで新しいモデルが出るので偶にチェックするのがおすすめ。\n使用できるモデルは\ncurl -s https://api.openai.com/v1/models -H \"Authorization: Bearer $OPENAI_API_KEY\" | gojq -r \".data[].id\" | grep gpt | sort\nで得られる。\n例) ai-model: gpt-3.5-turbo-16k\n(基本的に設定変更不要)"`
	Temperature float64 `json:"temperature" comment:"これは、AIに与える変数で、返答のランダム性を制御するパラメータである。\n値が小さいほどよくある解答など決まりきった解答を返し、\n値が大きいほど奇抜な返答が返ってきやすくなる。\n0から2の範囲の値を設定する。詳細は\nhttps://platform.openai.com/docs/api-reference/completions/create#completions/create-temperature\nを参照。"`
	OpenaiURL   string  `json:"openai-url"  comment:"OpenAIのAPIで使用するURL。\n基本的に\"https://api.openai.com/v1/chat/completions\"だが、\n将来変更されるかもしれないので設定できるようにしておく。\n例) openai-url: https://api.openai.com/v1/chat/completions\n(基本的に設定変更不要)"`
}
