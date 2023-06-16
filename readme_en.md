# easygpt

An app that uses the chatgpt API to process text files in bulk.  
It can be used for translation, adding comments to source code, and more, depending on how it is used.

## Install

```sh
go install github.com/xcd0/easygpt@latest
```

## Simple Example (Getting Output by Providing a String to the AI)

Set your OpenAI API key in the `OPENAI_API_KEY` environment variable.  
The API key below is an example and is invalid, so you need to generate your own valid API key and set it.  
You can generate an API key from https://platform.openai.com/account/api-keys.  
Set the generated API key in `.bashrc` or similar.
```sh
$ export OPENAI_API_KEY=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
$ echo export OPENAI_API_KEY=$OPENAI_API_KEY >> .bashrc
```

After that, install easygpt and run it.

```sh
$ go install github.com/xcd0/easygpt@latest
$ easygpt -i "Say this is test."
this is test.
```

## Bulk Processing of Multiple Files

Process multiple text files with AI in bulk.

### Preparation

Generate a template configuration file and edit it to create a configuration file.

1. Generate the template.

```sh
$ easygpt --create-setting
```

This will generate a `eagygpt.hjson` file in the current directory.  
You can also use `easygpt -c`.

2. Edit the configuration file.
* Setting the API key
	* It is recommended to set it in the `OPENAI_API_KEY` environment variable. In this case, there is no need to write it in the configuration file.
	* You can also write the API key in the `apikey` section of the configuration file.  
		Example: `apikey: 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`  
		The above API key is an invalid API key, so you need to generate your own.  
		You can generate an API key from https://platform.openai.com/account/api-keys.

* Write the string you want to give before the input text file in the `prompt` section of the configuration file.  
	For example, if you want to translate an English text file, write `prompt: Translate the following.`  
	If you want to write multiple lines, use `\n` to line break or use the following heredoc format.  
	```
	prompt:  
	'''  
	You can write multiple lines like this.  
	'''  
	```

Read the comments in the configuration file and write the necessary settings accordingly.

3. Configure the file.

The configuration file is searched according to the following specifications.  
Basically, just place the easygpt.hjson file in the current directory.

The default names for the configuration file are:  
* easygpt.hjson  
* .easygpt.hjson  
* .easygpt  
You can use a configuration file with a name other than these three by setting it with the `--setting` argument.

The default locations for the configuration file are:  
* Current directory  
* Home directory  
* The same directory as the executable file  
The configuration file is searched in this order.  
Therefore, if it exists in the current directory, the configuration file in the home directory will be ignored.

## Usage 1: Drag and Drop

Usage on a GUI. Select files or folders and drag and drop them onto the executable file.  
All the files in the given files and folders will be processed by the AI in bulk.  
The output will be saved with the input file name followed by `_easygpt_output`.

If a configuration file exists, you can use it as follows:
1. Select the text file or directory you want to submit to the AI and drag and drop it onto the easygpt executable file.
1. The processing result will be output with the input file name followed by `_easygpt_output`.

When using this method, the following settings in the configuration file will be ignored:
* input-dir
* output-dir
* extension

## Usage 2: Execute from Command Line

As a simple example of how to run, if the API key is set in the environment variable or configuration file, you can use it as follows:
```sh
$ ./easygpt -i Introduce yourself.
Hello, I am an AI. I am a natural language processing model developed by OpenAI. My purpose is to provide the best answers and responses when users ask questions or make requests. I have information on various topics and can also make grammar and style corrections. How can I assist you?
```

## Usage 3: Execute from Command Line and Process a Set of Files

If command line arguments are set, they will override the values in the configuration file.  
For values that are set in the configuration file in the default path, the arguments can be omitted.  
Since there are many arguments that can be set, it is recommended to use a configuration file in most cases.

There are flag arguments in the form of `--key value` and non-flag arguments that are not like that.  
* Flag arguments have priorities, and if a specific argument is used, other arguments may be ignored.  
  For example, if `--create-setting` is included, other arguments will be ignored and the program will only generate the template configuration file and then exit.  
  The value does not need to be enclosed in "", except when you want to leave it empty.
* Non-flag arguments are treated as input files.  
  If a directory is specified, all files in that directory will be treated as input files recursively.

The order of the arguments does not matter.  
For more information about the arguments, refer to the help text output by `./easygpt --help`.

## Open Source Libraries Used

The license for each open source library is kept in `./licenses`.

* Argument parsing and help generation
	* [alexflint/go-arg](https://github.com/alexflint/go-arg)
		* BSD-2-Clause license
* Random string generation (UUID is not required)
	* [google/uuid](https://github.com/google/uuid)
		* BSD-3-Clause license
* Stack trace
	* [pkg/errors](https://github.com/pkg/errors)
		* BSD 2-Clause "Simplified" License
* Configuration file
	* [hjson/hjson-go](https://github.com/hjson/hjson-go)
		* MIT License
* Go programming language
	* [golang/go](https://github.com/golang/go)
		* BSD 3-Clause "New" or "Revised" License
* License file retrieval
	* [google/go-licenses](https://github.com/google/go-licenses)
		* Apache License 2.0
* Executable file compression
	* [upx/upx](https://github.com/upx/upx)
		* GPL2+ or UPX LICENSE
			* https://upx.github.io/upx-license.html
		* The binary is used during the make process, and the resulting binary is distributed without modification. The program does not include the source code for upx.

## LICENSE

MIT License
