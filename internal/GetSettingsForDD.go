package internal

type ArgsDD struct {
	InputFiles []string `arg:"positional"         help:"非フラグ引数はファイルやディレクトリのパスであると見なす。\n                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、\n                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。"`
}

func GetSettingsForDD(argsAll *ArgsAll, settings []Setting) ([]SettingForDD, []string) {
	var setting Setting
	if len(settings) != 0 {
		setting = settings[0]
	} else {
		return nil, nil
	}

	// D&Dで実行された場合の対応のために、
	// 引数に--xxxのような形式でない引数があれば、
	// それは個別に実行する。

	argsDD := ArgsDD{
		InputFiles: argsAll.InputFiles,
	}
	// go-argはポジショナル引数(--aaa bbb形式でない、第n引数のように使用される引数)を解析したとき、ポジショナル引数が存在しない場合エラー終了してしまうので、
	// ポジショナル引数がなくても終了しないよう、arg.osExitをオーバーライドする。
	var err error = nil
	//arg.osExit = func() {
	//	err = errors.Errorf("引数によって直接指定された入力ファイルがありませんでした。")
	//	log.Printf("%+v", err)
	//}
	//arg.MustParse(&argsDD)
	if err != nil {
		// ポジショナル引数がなかった。別によい。
		return []SettingForDD{}, nil
	} else if len(argsDD.InputFiles) > 0 {
		// ポジショナル引数があった。
		//log.Printf("InputFiles: %v", argsDD.InputFiles)
		settingsfordd, files := GenerateSettingForDD(&setting.Apikey, &setting.Prompt, &setting.Postfix, &argsDD.InputFiles)
		return settingsfordd, files
	} else {
		return []SettingForDD{}, nil // ここには来ないはず。ここに来る意味がよくわからない。
	}
}
