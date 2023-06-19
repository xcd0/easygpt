# easygpt

A program that uses the chatgpt API to process text files in bulk.  
It can be used for translation, adding comments to source code, and more, depending on how it is used.

## Install

```sh
go install github.com/xcd0/easygpt@latest
```

## Preparation

Set your OpenAI API key to the environment variable `OPENAI_API_KEY`.  
The API key below is an invalid API key for example purposes, so generate your own valid API key and set it.  
API keys can be generated from https://platform.openai.com/account/api-keys.  
Set the generated API key in `.bashrc` or similar.
```sh
$ export OPENAI_API_KEY=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
$ echo export OPENAI_API_KEY=$OPENAI_API_KEY >> .bashrc
```

## Example 1: Providing a string to the AI and getting the output

If the API key is specified in the environment variable, you can execute it directly from the command line.

```sh
$ easygpt -i "Say this is test."
this is test.
$ easygpt -i 自己紹介してください。
こんにちは、私はAIです。私はOpenAIが開発した自然言語処理モデルです。私の目的は、ユーザーが質問や要求をすると、最善の回答や応答を提供することです。私は様々なトピックについての情報を持っており、文法やスタイルの修正も行うことができます。どのようにお手伝いできますか？
```

## Example 2: D&D (Drag and Drop)

Process all the files included in the given file or folder in bulk using the AI.  
By selecting files or folders and dragging and dropping them onto the executable file, you can process the included files as text with the AI.  
* Intended for use in a GUI environment.  
* The output is saved with the input file name followed by `_easygpt_output`.  
* The behavior can be modified by the configuration file described below.
	* For this usage, the following settings in the configuration file are ignored:
		* input-dir
		* output-dir
		* extension

## Example 3: Bulk processing of multiple files

Process multiple text files in a specified directory in bulk using the AI.  
This is the main usage.  
This usage specifies all the behavior in the configuration file.  
More detailed control can be achieved compared to Example 2.

## About the Configuration File

Generate a template for the configuration file and edit it to create the configuration file.

| key          | Required | Description                  |
|--------------|----------|------------------------------|
| input-dir    | Required | Input directory.<br>A directory that stores the files you want to give to the AI. |
| output-dir   | Required | Output directory.<br>A directory that will be created to have the same directory structure as the input directory. |
| apikey       | Optional | Specify the API key. Not required if specified in the environment variable.<br>If specified in the configuration file, it will override the environment variable. |
| prompt       | Optional | A string specified here will be prepended to the beginning of all input files.<br>Can be used to provide instructions to the AI, etc. |
| extension    | Optional | Limit the input files in the input directory to the specified extension.<br>If this is empty, "*", or not specified, there is no restriction based on the extension. |
| concurrency  | Optional | The number of concurrent executions. Be aware of the API rate limit. |
| temperature  | Optional | Specify the AI's `temperature` variable. |
| move         | Optional | Move successfully processed files to the specified directory.<br>This makes it easier to resume execution from the middle in case of errors or interruptions such as Ctrl-C. |
| tmp-dir      | Optional | Generally not necessary to specify. A directory to store temporary files. |
| postfix      | Optional | Generally not necessary to specify. Specify a string to append to the end of the output file name. |
| ai-model     | Optional | Generally not necessary to specify. Specify the AI model you want to use. |
| openai-url   | Optional | Generally not necessary to specify. Specify the API URL. |


1. Generating a template.

Since there are many configurable items, it is recommended to generate a template for the configuration file and edit it.

The following command generates `eagygpt.hjson` in the current directory.
```sh
$ easygpt --create-setting
```
`easygpt -c` is also acceptable.


2. Editing the configuration file.
* Setting the API key
	* You can specify the API key in the configuration file. Not necessary if set in the environment variable.

* In the `prompt` section of the configuration file, write the string you want to prepend to the input text files.  
   For example, if you want to translate an English text file, write `prompt: Please translate the following.`  
   When you want to write multiple lines, use `\n` to line break or use the following heredoc format.  
   ```
   prompt:  
   '''  
   You can write multiple lines like this.  
   '''  
   ```
   For more details, refer to the `Mulutiline Setting` section in the [hjson documentation](https://hjson.github.io/syntax.html).

For other arguments, read the comments in the generated configuration file and write as necessary.  
Basically, `input-dir` and `output-dir` should be specified.

3. Placing the configuration file.

Basically, you just need to place `easygpt.hjson` in the current directory.

The configuration file is searched according to the following specifications.

The default names for the configuration file are:
* easygpt.hjson
* .easygpt.hjson
* .easygpt

If you want to use a configuration file with a different name, you can specify it with the `--setting` argument.

The default locations for the configuration file are:
* Current directory
* Home directory
* The same directory as the executable file

The configuration file is searched in this order.
Therefore, if it exists in the current directory, the configuration file in the home directory will be ignored.


## Configuration via Command Line Arguments

If command line arguments are set, the values in the configuration file will be overwritten with the values in the command line arguments.  
For values that are set in the configuration file in the default path, the arguments can be omitted.  
Since there are many arguments that can be set, it is recommended to use the configuration file in most cases.

There are flag arguments in the form of `--key value` and non-flag arguments that are not like that.  
* Flag arguments have priority, and if a specific argument is used, other arguments may be ignored.  
  For example, if `--create-setting` is present, it will ignore other arguments and generate the template for the configuration file and exit.  
  Unless you want to set the value to empty, you don't need to enclose it in "".  
* Non-flag arguments are treated as input files.  
  If a directory is specified, all files in that directory will be treated as input files recursively.

The order of the arguments does not matter.  
For more details about the arguments, refer to the help text output by `./easygpt --help`.


## OSS Libraries Used

The licenses of each OSS library are stored in `./licenses`.

* Argument parsing and help generation
	* [alexflint/go-arg](https://github.com/alexflint/go-arg)
		* BSD-2-Clause license
* Random string generation (does not have to be UUID specifically)
	* [google/uuid](https://github.com/google/uuid)
		* BSD-3-Clause license
* Stack trace
	* [pkg/errors](https://github.com/pkg/errors)
		* BSD 2-Clause "Simplified" License
* Configuration file
	* [hjson/hjson-go](https://github.com/hjson/hjson-go)
		* MIT License
* Go language itself
	* [golang/go](https://github.com/golang/go)
		* BSD 3-Clause "New" or "Revised" License
* Obtaining license files
	* [google/go-licenses](https://github.com/google/go-licenses)
		* Apache License 2.0
* Executable file compression
	* [upx/upx](https://github.com/upx/upx)
		* GPL2+ or UPX LICENSE
			* https://upx.github.io/upx-license.html
		* The program is distributed without modification by compressing it with upx during the build process. The source code of upx is not included in the program.

## LICENSE

MIT License