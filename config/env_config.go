package config
/*
	read file(.env) var to set env(linux) var
*/

type EnvOptions struct{
	FileType string		//file type. such as yaml/json/env etc.
	FilePath string		//file path. runtime directory as root.
	EnvPrefix string	//env variable pre. such as preXXX
}