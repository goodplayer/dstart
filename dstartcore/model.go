package dstartcore

type Config struct {
	Global Global `toml:"global"`

	Services map[string]Service `toml:"service"`
}

type Global struct {
	IncludeConfigFileExp string            `toml:"IncludeConfigFileExp"`
	Envs                 map[string]string `toml:"envs"`
}

type Service struct {
	Description string `toml:"Description"`
	ExecCommand string `toml:"ExecCommand"`

	Username  string `toml:"Username"`
	Uid       string `toml:"Uid"`
	GroupName string `toml:"Groupname"`
	GroupId   string `toml:"GroupId"`

	WorkingDirectory string `toml:"WorkingDirectory"`

	InputFile                   string `toml:"InputFile"`
	OutputFile                  string `toml:"OutputFile"`
	OutputFileAppend            bool   `toml:"OutputFileAppend"`
	ErrorOutputFile             string `toml:"ErrorOutputFile"`
	ErrorOutputFileAppend       bool   `toml:"ErrorOutputFileAppend"`
	RedirectErrorOutputToOutput bool   `toml:"RedirectErrorOutputToOutput"`

	Envs map[string]string `toml:"Envs"`
}
