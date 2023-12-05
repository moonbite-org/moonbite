package config

type config struct {
	ArchiveName  string
	VersionStamp string
}

var Config = config{
	ArchiveName:  "<moonbite archive>",
	VersionStamp: "0.0.1-pre-alpha",
}
